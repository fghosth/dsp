package rabbitmq

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/streadway/amqp"
	"jvole.com/dsp/util"
)

var logger *log.Logger

func init() {
	logger = &util.KitLogger
	log.With(*logger, "component", "RabbitMQ")
}

// Receiver 观察者模式需要的接口
// 观察者用于接收指定的queue到来的数据
type Receiver interface {
	QueueName() string          // 获取接收者需要监听的队列
	RouterKey() string          // 这个队列绑定的路由
	OnError(err error)          // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
	OnReceive(body []byte) bool // 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
}

// RabbitMQ 用于管理和维护rabbitmq的对象
type RabbitMQ struct {
	wg           sync.WaitGroup
	source       *amqp.Connection
	channel      *amqp.Channel
	exchangeName string // exchange的名称
	exchangeType string // exchange的类型
	receivers    []Receiver
	server       Mqserver
}
type Mqserver struct {
	Url    string
	Online bool //是否在线
}

var RabbitMQConn *RabbitMQ

func (mq *RabbitMQ) Publish(exname, routekey string, body []byte) (success bool) {

	//发送消息
	err := mq.channel.Publish(
		exname,   // exchange     交换器名称，使用默认
		routekey, // routing key    路由键，这里为队列名称
		false,    // mandatory
		false,
		amqp.Publishing{
			ContentType: "text/plain", //消息类型，文本消息
			Body:        body,
		})
	if err != nil {
		success = false
	} else {
		success = true
	}
	return
}

// New 创建一个新的操作RabbitMQ的对象
func NewRabbitMQ(exname, extype string, server Mqserver) *RabbitMQ {
	rbmq := &RabbitMQ{
		exchangeName: exname,
		exchangeType: extype,
		server:       server,
	}
	rbmq.receivers = make([]Receiver, 0)
	conn, err := amqp.Dial(rbmq.server.Url)

	if err != nil {
		return nil
	} else {
		rbmq.source = conn
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil
	} else {
		rbmq.channel = ch
	}

	// 这里可以根据自己的需要去定义
	return rbmq
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
func (mq *RabbitMQ) RegisterReceiver(receiver Receiver) {

	mq.receivers = append(mq.receivers, receiver)
}

// prepareExchange 准备rabbitmq的Exchange
func (mq *RabbitMQ) prepareExchange() error {

	// 申明Exchange
	err := mq.channel.ExchangeDeclare(
		mq.exchangeName, // exchange
		mq.exchangeType, // type
		true,            // durable
		false,           // autoDelete
		false,           // internal
		false,           // noWait
		nil,             // args
	)

	if nil != err {

		return err
	}
	return nil
}

//检查连接是否正常
func (mq *RabbitMQ) Refresh() bool {
	_, err := amqp.Dial(mq.server.Url)
	if err != nil {
		return false
	} else {
		return true
	}
}

//重新连接channel
func (mq *RabbitMQ) Channel() *amqp.Channel {

	conn, err := amqp.Dial(mq.server.Url)
	if err != nil {
		return nil
	} else {
		mq.source = conn
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil
	} else {
		return ch
	}
}

//销毁连接对象
func (mq *RabbitMQ) Distory() {
	mq.Channel().Close()
	mq.source.Close()
}

// run 开始获取连接并初始化相关操作
func (mq *RabbitMQ) run() {
	if !mq.Refresh() {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "rabbit刷新连接失败，将要重连",
		)
		return
	}

	// 获取新的channel对象
	mq.channel = mq.Channel()
	// pp.Println(mq.channel)
	// 初始化Exchange
	mq.prepareExchange()

	for _, receiver := range mq.receivers {
		mq.wg.Add(1)

		go mq.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
	}

	mq.wg.Wait()

	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	level.Error(*logger).Log(
		"method", f.Name(),
		"file", path.Base(file),
		"line", line,
		"msg", "所有处理queue的任务都意外退出了",
	)
	// 理论上mq.run()在程序的执行过程中是不会结束的
	// 一旦结束就说明所有的接收者都退出了，那么意味着程序与rabbitmq的连接断开
	// 那么则需要重新连接，这里尝试销毁当前连接
	defer mq.Distory()
}

// Start 启动Rabbitmq的客户端
func (mq *RabbitMQ) Start() {
	for {

		mq.run()

		// 一旦连接断开，那么需要隔一段时间去重连
		time.Sleep(3 * time.Second)
	}
}

// Listen 监听指定路由发来的消息
// 这里需要针对每一个接收者启动一个goroutine来执行listen
// 该方法负责从每一个接收者监听的队列中获取数据，并负责重试
func (mq *RabbitMQ) listen(receiver Receiver) {
	defer mq.wg.Done()

	// 这里获取每个接收者需要监听的队列和路由
	queueName := receiver.QueueName()
	routerKey := receiver.RouterKey()
	// 申明Queue
	_, err := mq.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive(排他性队列)
		false,     // no-wait
		nil,       // arguments
	)
	if nil != err {
		// 当队列初始化失败的时候，需要告诉这个接收者相应的错误
		receiver.OnError(fmt.Errorf("初始化队列 %s 失败: %s", queueName, err.Error()))
	}

	// 将Queue绑定到Exchange上去
	err = mq.channel.QueueBind(
		queueName,       // queue name
		routerKey,       // routing key
		mq.exchangeName, // exchange
		false,           // no-wait
		nil,
	)

	if nil != err {
		receiver.OnError(fmt.Errorf("绑定队列 [%s - %s] 到Exchanges失败: %s", queueName, routerKey, err.Error()))
	}

	// 获取消费通道
	//prefetchSize：0 prefetchSize maximum amount of content (measured in* octets) that the server will deliver, 0 if unlimited
	//prefetchCount：会告诉RabbitMQ不要同时给一个消费者推送多于N个消息，即一旦有N个消息还没有ack，则该consumer将block掉，直到有消息ack
	//global：true\false 是否将上面设置应用于channel，简单点说，就是上面限制是channel级别的还是consumer级别
	//备注：据说prefetchSize 和global这两项，rabbitmq没有实现，暂且不研究
	mq.channel.Qos(10, 0, false) // 确保rabbitmq会一个一个发消息
	msgs, err := mq.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("获取队列 %s 的消费通道失败: %s", queueName, err.Error()))
	}

	// 使用callback消费数据
	for msg := range msgs {
		// 当接收者消息处理失败的时候，
		// 比如网络问题导致的数据库连接失败，连接失败等等这种
		// 通过重试可以成功的操作，那么这个时候是需要重试的
		// 直到数据处理成功后再返回，然后才会回复rabbitmq ack
		for !receiver.OnReceive(msg.Body) {
			pc, file, line, _ := runtime.Caller(1)
			f := runtime.FuncForPC(pc)
			level.Error(*logger).Log(
				"method", f.Name(),
				"file", path.Base(file),
				"line", line,
				"msg", "receiver 数据处理失败，将要重试",
			)
			time.Sleep(1 * time.Second)
		}

		// 确认收到本条消息, multiple必须为false
		msg.Ack(false)
	}
}
