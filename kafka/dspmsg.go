package kafka

import (
	"log"

	"github.com/Shopify/sarama"
	"jvole.com/dsp/config"
)

type dspMsg struct{}

func NewDspMsg() *producer {
	//设置配置
	cfg := sarama.NewConfig()

	//等待服务器所有副本都保存成功后的响应
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	// cfg.Producer.Partitioner = sarama.NewHashPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	cfg.Version = sarama.V0_11_0_0
	producer := NewProducer([]string{config.KafkaURL}, cfg, success, errfun)
	return producer
}

func success(suc *sarama.ProducerMessage) {
	log.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
}

func errfun(err *sarama.ProducerError) {
	log.Printf(err.Error(), err.Msg)
}
