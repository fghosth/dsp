package kafka

import (
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var KafKaProducer *producer

type producer struct {
	/*
		c.Net.MaxOpenRequests = 5
		c.Net.DialTimeout = 30 * time.Second
		c.Net.ReadTimeout = 30 * time.Second
		c.Net.WriteTimeout = 30 * time.Second
		c.Net.SASL.Handshake = true

		c.Metadata.Retry.Max = 3
		c.Metadata.Retry.Backoff = 250 * time.Millisecond
		c.Metadata.RefreshFrequency = 10 * time.Minute
		c.Metadata.Full = true

		c.Producer.MaxMessageBytes = 1000000
		c.Producer.RequiredAcks = WaitForLocal
		c.Producer.Timeout = 10 * time.Second
		c.Producer.Partitioner = NewHashPartitioner  //选择分区的分区选择器.用于选择主题的分区
		c.Producer.Retry.Max = 3 //重试次数
		c.Producer.Retry.Backoff = 100 * time.Millisecond
		c.Producer.Return.Errors = true  //是否接收返回的错误消息,当发生错误时会放到Error这个通道中.从它里面获取错误消息

		//抓取数据的大小设置
		c.Consumer.Fetch.Min = 1
		c.Consumer.Fetch.Default = 32768

		c.Consumer.Retry.Backoff = 2 * time.Second //失败后再次尝试的间隔时间
		c.Consumer.MaxWaitTime = 250 * time.Millisecond  //最大等待时间
		c.Consumer.MaxProcessingTime = 100 * time.Millisecond
		c.Consumer.Return.Errors = false  //是否接收返回的错误消息,当发生错误时会放到Error这个通道中.从它里面获取错误消息
		c.Consumer.Offsets.CommitInterval = 1 * time.Second // 提交跟新Offset的频率
		c.Consumer.Offsets.Initial = OffsetNewest // 指定Offset,也就是从哪里获取消息,默认时从主题的开始获取.

		c.ClientID = defaultClientID
		c.ChannelBufferSize = 256  //通道缓存大小
		c.Version = minVersion //指定kafka版本,不指定,使用最小版本,高版本的新功能可能无法正常使用.
		c.MetricRegistry = metrics.NewRegistry()
	*/
	successFun func(*sarama.ProducerMessage) //异步模式时需要传的参数
	errorFun   func(*sarama.ProducerError)   //异步模式时需要传的参数
	prodAsync  sarama.AsyncProducer
	prodSync   sarama.SyncProducer
}

func NewProducer(server []string, config *sarama.Config, successfun func(*sarama.ProducerMessage), errfun func(*sarama.ProducerError)) *producer {
	prodAsync, e := sarama.NewAsyncProducer(server, config)
	if e != nil {
		panic(e)
	}
	prodSync, e := sarama.NewSyncProducer(server, config)
	if e != nil {
		panic(e)
	}
	pd := &producer{
		successfun,
		errfun,
		prodAsync,
		prodSync,
	}
	go pd.run()
	return pd
}

//发送消息
func (pd *producer) SendMsgSync(message, key, topic string, t time.Time) error {
	//发送的消息,主题,key
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(message),
		Timestamp: t,
	}
	_, _, err := pd.prodSync.SendMessage(msg)
	// part, offset, err := pd.prodSync.SendMessage(msg)
	if err != nil {
		return err
		// log.Printf("send message(%s) err=%s \n", message, err)
	} else {
		// log.Printf("发送成功，partition=%d, offset=%d \n", part, offset)
		return nil
	}
}

//发送消息
func (pd *producer) SendMsgAsync(message, key, topic string, t time.Time) {
	//发送的消息,主题,key
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(message),
		Timestamp: t,
	}
	//使用通道发送
	pd.prodAsync.Input() <- msg
}

func (pd *producer) run() {
	var wg sync.WaitGroup
	wg.Add(2) //2 goroutine
	// 发送成功
	go func() {
		defer wg.Done()
		for v := range pd.prodAsync.Successes() {
			pd.successFun(v)
		}
	}()

	// 发送失败
	go func() {
		defer wg.Done()
		for err := range pd.prodAsync.Errors() {
			pd.errorFun(err)
		}
	}()

	wg.Wait()
}

//关闭连接
func (pd *producer) Close() {
	pd.prodAsync.AsyncClose()
	pd.prodSync.Close()
}
