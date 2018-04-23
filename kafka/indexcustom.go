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

type indexReceiver struct {
	topic []string
	// partition int32
	// offset    int64
	groupID    string
	offsets    []OffsetGroup
	partition  mapset.Set
	timeFilter map[string]time.Time
}

const NUM_TimeFilter = 10

//接受的消息结构
type msgRec struct {
	Action string `json:"action"`
	CID    uint32 `json:"cid"`
	T      int64  `json:"t"`
	ADX    uint32 `json:"adx"`
}

func NewIndexReceiver(topic []string, groupID string) *indexReceiver {
	return &indexReceiver{
		topic:     topic,
		groupID:   groupID,
		offsets:   make([]OffsetGroup, 0),
		partition: mapset.NewSet(),
	}
}

// 设置offset
func (ic *indexReceiver) SetOffset(offsets []OffsetGroup) {
	ic.offsets = offsets
}

// 设置offset
func (ic *indexReceiver) SetPartition(part []int32) {
	for _, v := range part {
		ic.partition.Add(v)
	}
}

//获取时间过滤条件，最多十个
func (ic *indexReceiver) TimeFilter() map[string]time.Time {
	return ic.timeFilter
}

//设置时间过滤条件，最多十个
func (ic *indexReceiver) SetTimeFilter(tf map[string]time.Time) (err error) {
	if len(tf) > NUM_TimeFilter { //超过10个返回错误
		err = errors.New("timeFilter max len is " + strconv.Itoa(NUM_TimeFilter))
		return
	}
	ic.timeFilter = tf
	return
}

// 获取接收者需要监听的队列topic
func (ic *indexReceiver) Topic() []string {
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
func (ic *indexReceiver) GroupID() string {
	return ic.groupID
}

// 这个队列offset
func (ic *indexReceiver) Offsets() []OffsetGroup {

	return ic.offsets
}

// 这个队列offset
func (ic *indexReceiver) Partition() mapset.Set {

	return ic.partition
}

// 处理遇到的错误，当kafka对象发生了错误，他需要告诉接收者处理错误
func (ic *indexReceiver) OnError(err error) {
	log.Println("indexconsumer", err)
}

// 处理收到的消息, 这里需要告知kafka对象消息是否处理成功
func (ic *indexReceiver) OnReceive(msg *sarama.ConsumerMessage) (success bool) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	fmt.Println("msg offset: ", msg.Offset, " partition: ", msg.Partition, " timestrap: ", msg.Timestamp.Format("2006-Jan-02 15:04"), " value: ", string(msg.Value))
	success = true
	return
}
