package kafka

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/deckarep/golang-set"
)

// Receiver 观察者模式需要的接口
// 观察者用于接收指定的queue到来的数据
type Receiver interface {
	Topic() []string                        // 获取接收者需要监听的topic
	GroupID() string                        //offset
	Offsets() []OffsetGroup                 //不同group，topic的offset信息
	Partition() mapset.Set                  //接受的partition
	TimeFilter() map[string]time.Time       //时间范围，string只有『After』[Befter](不包含，开区间)，支持十个
	OnError(err error)                      // 处理遇到的错误，当kafka对象发生了错误，他需要告诉接收者处理错误
	OnReceive(*sarama.ConsumerMessage) bool // 处理收到的消息, 这里需要告知kafka对象消息是否处理成功
}

type consumer struct {
	wg        *sync.WaitGroup
	receivers []Receiver
	server    []string
	config    *cluster.Config
}

type OffsetGroup struct {
	Topic     string
	Partition int32
	Offset    int64
	Metadata  string
}

func NewConsumer(server []string, config *cluster.Config) *consumer {
	config.Group.Mode = cluster.ConsumerModePartitions
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	kafka := &consumer{
		new(sync.WaitGroup),
		make([]Receiver, 0),
		server,
		config,
	}
	return kafka
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
func (cs *consumer) RegisterReceiver(receiver Receiver) {
	cs.receivers = append(cs.receivers, receiver)
}

//销毁连接对象
// func (cs *consumer) Close() {
// 	cs.cons.Close()
// }

// run 开始获取连接并初始化相关操作
func (cs *consumer) run() {
	for _, receiver := range cs.receivers {
		cs.wg.Add(1)
		go cs.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
	}

	cs.wg.Wait()
	// defer cs.Close()
}

// Start 启动kafka的客户端
func (cs *consumer) Start() {
	for {

		cs.run()

		// 一旦连接断开，那么需要隔一段时间去重连
		time.Sleep(3 * time.Second)
	}
}

func (cs *consumer) setOffset(offsets []OffsetGroup, groupid string, topics []string) (err error) {
	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Group.Mode = cluster.ConsumerModeMultiplex
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Version = sarama.V0_11_0_0
	c, err := cluster.NewConsumer(cs.server, groupid, topics, config)
	defer c.Close()
	if err != nil {
		// log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupid, err)
		return err
	}
	// consume errors
	go func() {
		for err := range c.Errors() {
			log.Printf("setoffset---%s:Error: %s\n", groupid, err.Error())
		}
	}()
	// consume notifications
	go func() {
		for ntf := range c.Notifications() {
			log.Printf("setoffset---%s:Rebalanced: %+v \n", groupid, ntf)
		}
	}()
	// for _, v := range offsets {
	// 	c.ResetPartitionOffset(v.Topic, v.Partition, v.Offset, v.Metadata)
	// }

	_, ok := <-c.Messages()
	if ok {
		for _, v := range offsets {
			c.ResetPartitionOffset(v.Topic, v.Partition, v.Offset, v.Metadata)
		}
	} else {
		err = errors.New("setoffset error")
	}
	return
}

// Listen 监听指定topic发来的消息
// 该方法负责从每一个接收者监听的队列中获取数据
func (cs *consumer) listen(receiver Receiver) {
	defer cs.wg.Done()
	// 这里获取每个接收者需要监听的topic,Partition,Offset
	topics := receiver.Topic()
	// partition := receiver.Partition()
	// offset := receiver.Offset()
	groupid := receiver.GroupID()
	offsets := receiver.Offsets()
	partition := receiver.Partition()
	timeFilter := receiver.TimeFilter()
	if len(offsets) > 0 { //如果需要设置offset
		cs.setOffset(offsets, groupid, topics)
	}
	cs.config.Group.Mode = cluster.ConsumerModePartitions
	//根据消费者获取指定的主题分区的消费者,Offset这里指定为获取最新的消息.
	c, err := cluster.NewConsumer(cs.server, groupid, topics, cs.config)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupid, err)
		return
	}
	defer c.Close()

	// consume errors
	go func() {
		for err := range c.Errors() {
			receiver.OnError(err)
			// log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()
	// consume notifications
	go func() {
		for ntf := range c.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupid, ntf)
		}
	}()

	//循环等待接受消息.
	for {
		select {
		//接收消息通道和错误通道的内容.
		case part, ok := <-c.Partitions():
			if ok {
				// pp.Println(part.Partition(), partition.Cardinality())
				//start a separate goroutine to consume messages
				if !partition.Contains(part.Partition()) && partition.Cardinality() > 0 {
					continue
				}
				go func(pc cluster.PartitionConsumer) {
					for msg := range pc.Messages() {
						//过滤时间条件
						if !cs.timeFilter(timeFilter, msg.Timestamp) && len(timeFilter) > 0 { //如果设置了timefilter，并且过滤失败不返回此消息
							continue
						}
						receiver.OnReceive(msg)
						// fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
						c.MarkOffset(msg, "") // mark message as processed
					}
				}(part)
			}

		}
	}
}

//判断t是否tf条件范围内
func (cs *consumer) timeFilter(tf map[string]time.Time, t time.Time) (res bool) {
	res = true
	for k, v := range tf {
		switch k {
		case "After":
			if !t.After(v) {
				res = false
				return
			}
		case "Before":
			if !t.Before(v) {
				res = false
				return
			}
		}
	}
	return
}
