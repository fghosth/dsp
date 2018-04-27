package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/config"
	"jvole.com/dsp/index"
)

type successReceiver struct {
	queueName string
	routeKey  string
}

//接受的消息结构
type WinNotify struct {
	ID      string `json:"id"`      //广告位id
	Price   uint32 `json:"price"`   //竞价价格
	CID     uint32 `json:"cid"`     //compaignid
	UID     uint32 `json:"uid"`     //用户id
	Postion string `json:"postion"` //广告位置标识
	Device  string `json:"device"`  //设备标识
	T       int64  `json:"t"`       //竞价成功时间
	User    string `json:"user"`    //广告端user标识
	ADX     uint32 `json:"adx"`     //adx
	ClickID string `json:"clickID"` //clickid
}

func NewSuccessReceiver(queueName, routeKey string) *successReceiver {
	return &successReceiver{queueName, routeKey}
}

// 获取接收者需要监听的队列
func (ic successReceiver) QueueName() string {
	return ic.queueName
}

// 这个队列绑定的路由
func (ic successReceiver) RouterKey() string {
	return ic.routeKey
}

// 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
func (ic successReceiver) OnError(err error) {
	log.Println(err)
}

// 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
func (ic successReceiver) OnReceive(body []byte) (success bool) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	msg := &WinNotify{}
	err := json.Unmarshal(body, &msg)
	if err != nil { //解析消息错误丢弃消息
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Debug(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"err", err,
			"msgbody", string(body),
			"msg", "successCustom解析消息错误",
		)
		return true
	}
	if config.StartupTime.Unix() > msg.T || (msg.ADX&config.ADXCode[config.ADX]) == 0 { //如果消息时间早于系统启动时间，消息直接丢弃
		// pp.Println("时间过期或不是此adx的compaign", msg.ADX, config.ADX)
		return true
	}
	if (config.ADXCode[config.ADX] & msg.ADX) == 0 { //不属于此adx
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Debug(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"compaignid", msg.CID,
			"msgbody", string(body),
			"msg", "次compaign不属于"+config.ADX,
		)
		return true
	}
	if msg.ID == "" { //消息无效丢弃消息
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Debug(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"err", err,
			"msgbody", string(body),
			"msg", "消息无效",
		)
		return true
	}
	cp := index.CPINDEX.GetCompaign(msg.CID)
	if cp.ID == 0 { //当缓存没有时 丢弃消息
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Debug(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"compaignid", msg.CID,
			"msgbody", string(body),
			"msg", "索引中没有此compaign",
		)
		return true
	}
	// pp.Println(string(body))
	//更新索引
	err = index.CPINDEX.BidSuccess(msg.CID, msg.Price, msg.Postion, msg.User, msg.Device)
	if err != nil {
		success = false
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "更新缓存索引失败",
		)
	} else {
		success = true
	}
	return
}
