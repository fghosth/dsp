package main

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"jvole.com/dsp/kafka"
)

func main() {
	topics := []string{"testTime"}
	groupID := "derek"
	receiver := kafka.NewIndexReceiver(topics, groupID)

	// //设置offset，默认不用设置，上次读取到的offset
	// offset := make([]kafka.OffsetGroup, 0)
	// os := *&kafka.OffsetGroup{
	// 	Topic:     topics[0],
	// 	Partition: 1,
	// 	Offset:    599,
	// 	Metadata:  "",
	// }
	// offset = append(offset, os)
	// receiver.SetOffset(offset)
	//接受哪些partition，默认不用设置，全部接受
	partition := []int32{0, 1, 2}
	receiver.SetPartition(partition)
	// //设置接受时间范围 默认不用设置，全部接受
	// tf := make(map[string]time.Time)
	// ti, _ := time.Parse("2006-01-02 15:04:05 -0700", "2018-04-22 21:30:00 +0800")
	// ti2, _ := time.Parse("2006-01-02 15:04:05 -0700", "2018-04-22 22:20:00 +0800")
	// tf["After"] = ti
	// tf["Before"] = ti2
	// receiver.SetTimeFilter(tf)
	var server = []string{"kafka.newbidder.com:9092"}
	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// config.Group.Mode = cluster.ConsumerModePartitions
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0
	consu := kafka.NewConsumer(server, config)
	consu.RegisterReceiver(receiver)
	consu.Start()
}
