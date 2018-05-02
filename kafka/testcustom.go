package kafka

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/deckarep/golang-set"
)

type testReceiver struct {
	topic []string
	// partition int32
	// offset    int64
	groupID    string
	offsets    []OffsetGroup
	partition  mapset.Set
	timeFilter map[string]time.Time
}

func NewTestReceiver(topic []string, groupID string) *testReceiver {
	return &testReceiver{
		topic:     topic,
		groupID:   groupID,
		offsets:   make([]OffsetGroup, 0),
		partition: mapset.NewSet(),
	}
}

// 设置offset
func (ic *testReceiver) SetOffset(offsets []OffsetGroup) {
	ic.offsets = offsets
}

// 设置offset
func (ic *testReceiver) SetPartition(part []int32) {
	for _, v := range part {
		ic.partition.Add(v)
	}
}

//获取时间过滤条件，最多十个
func (ic *testReceiver) TimeFilter() map[string]time.Time {
	return ic.timeFilter
}

//设置时间过滤条件，最多十个
func (ic *testReceiver) SetTimeFilter(tf map[string]time.Time) (err error) {
	if len(tf) > NUM_TimeFilter { //超过10个返回错误
		err = errors.New("timeFilter max len is " + strconv.Itoa(NUM_TimeFilter))
		return
	}
	ic.timeFilter = tf
	return
}

// 获取接收者需要监听的队列topic
func (ic *testReceiver) Topic() []string {
	return ic.topic
}

//// 这个队列绑定partition
// func (ic indexReceiver) Partition() int32 {
// 	return ic.partition
// }
//
// // 这个队列offset
// func (ic indexReceiver) Offset() int64 {
// 	return ic.offset
// }

// 这个队列groupid
func (ic *testReceiver) GroupID() string {
	return ic.groupID
}

// 这个队列offset
func (ic *testReceiver) Offsets() []OffsetGroup {

	return ic.offsets
}

// 这个队列offset
func (ic *testReceiver) Partition() mapset.Set {

	return ic.partition
}

// 处理遇到的错误，当kafka对象发生了错误，他需要告诉接收者处理错误
func (ic *testReceiver) OnError(err error) {
	log.Println("indexconsumer", err)
}

// 处理收到的消息, 这里需要告知kafka对象消息是否处理成功
func (ic *testReceiver) OnReceive(msg *sarama.ConsumerMessage) (success bool) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	// pp.Println(msg.Timestamp.Nanosecond(), time.Now().Nanosecond(), msg.Timestamp.Add(time.Nanosecond*460124).Nanosecond(), time.Now().Nanosecond()%1000000)
	// fmt.Println("key:", string(msg.Key), "msg offset: ", msg.Offset, " partition: ", msg.Partition, " timestrap: ", msg.Timestamp.Format("2006-Jan-02 15:04"), " value: ", string(msg.Value))

	tags := make(map[string]string)
	tags["adx"] = string(msg.Key)

	fields := make(map[string]interface{})
	fields["data"] = string(msg.Value)
	table := "testdata"
	times := msg.Timestamp.Add(time.Nanosecond * time.Duration(time.Now().Nanosecond()%1000000)) //补足纳秒时间
	err := InfluxConn.Insert(tags, fields, table, times)
	if err != nil {
		success = false
		log.Println(err)
	} else {
		success = true
	}

	return
}
