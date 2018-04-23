package kafka_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"jvole.com/dsp/kafka"
)

func aTestConsumer(t *testing.T) {
	//配置
	config := sarama.NewConfig()
	//接收失败通知
	config.Consumer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0
	//新建一个消费者
	consumer, e := sarama.NewConsumer([]string{"kafka.newbidder.com:9092"}, config)
	if e != nil {
		panic("error get consumer")
	}
	defer consumer.Close()

	//根据消费者获取指定的主题分区的消费者,Offset这里指定为获取最新的消息.
	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("error get partition consumer", err)
	}
	defer partitionConsumer.Close()
	//循环等待接受消息.
	for {
		select {
		//接收消息通道和错误通道的内容.
		case msg := <-partitionConsumer.Messages():
			fmt.Println("msg offset: ", msg.Offset, " partition: ", msg.Partition, " timestrap: ", msg.Timestamp.Format("2006-Jan-02 15:04"), " value: ", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Println(err.Err)
		}
	}
}
func aTestSyncSend(t *testing.T) {
	//设置配置
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//根据key的值取hash
	// config.Producer.Partitioner = sarama.NewHashPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0
	var server = []string{"kafka.newbidder.com:9092"}
	producer := kafka.NewProducer(server, config, nil, nil) //同步消息不用传最后2个参数
	count := 5
	for i := 0; i < count; i++ {
		message := "-------===-----" + strconv.Itoa(i)
		key := "11dspasdfew"
		topic := "test"
		t := time.Now()
		err := producer.SendMsgSync(message, key, topic, t)
		if err != nil {
			fmt.Println(err)
		}
	}

	producer.Close() //关闭发送器
}

func TestAsyncSend(t *testing.T) {
	//设置配置
	config := sarama.NewConfig()

	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	// config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Partitioner = sarama.NewHashPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0
	var server = []string{"kafka.newbidder.com:9092"}
	producer := kafka.NewProducer(server, config, success, errfun)
	count := 3
	for i := 0; i < count; i++ {
		message := "===5435345===" + strconv.Itoa(i)
		key := "dsp"
		topic := "testTime"
		// t := time.Now()
		t, _ := time.Parse("2006-01-02 15:04:05 -0700", "2018-04-22 21:33:00 +0400")
		producer.SendMsgAsync(message, key, topic, t)
	}

	time.Sleep(time.Second * 3)
	producer.Close() //关闭发送器
}

func success(suc *sarama.ProducerMessage) {
	fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
}

func errfun(err *sarama.ProducerError) {
	fmt.Printf(err.Error(), err.Msg)
}
