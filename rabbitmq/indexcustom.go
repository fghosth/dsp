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
	"jvole.com/dsp/model"
)

type indexReceiver struct {
	queueName string
	routeKey  string
}

//接受的消息结构
type msgRec struct {
	Action string `json:"action"`
	CID    uint32 `json:"cid"`
	T      int64  `json:"t"`
	ADX    uint32 `json:"adx"`
}

func NewIndexReceiver(queueName, routeKey string) *indexReceiver {

	return &indexReceiver{queueName, routeKey}
}

// 获取接收者需要监听的队列
func (ic indexReceiver) QueueName() string {
	return ic.queueName
}

// 这个队列绑定的路由
func (ic indexReceiver) RouterKey() string {
	return ic.routeKey
}

// 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
func (ic indexReceiver) OnError(err error) {
	log.Println(err)
}

// 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
func (ic indexReceiver) OnReceive(body []byte) (success bool) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()

	var msg msgRec
	err := json.Unmarshal(body, &msg)
	if err != nil { //解析消息错误丢弃消息
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"err", err,
			"msgbody", string(body),
			"msg", "解析消息错误",
		)
		return true
	}

	if config.StartupTime.Unix() > msg.T || (msg.ADX&config.ADXCode[config.ADX]) == 0 { //如果消息时间早于系统启动时间，或不属于此服务的消息直接丢弃
		// pp.Println("时间过期或不是此adx的compaign", msg.ADX, config.ADX)
		index.CPINDEX.Remove(msg.CID)
		return true
	}
	if (config.ADXCode[config.ADX] & msg.ADX) == 0 { //不属于此adx
		index.CPINDEX.Remove(msg.CID) //删除此adx
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
	switch msg.Action {
	case config.MQNADD, config.MQNMODIFY:
		cm := *&model.Compaign{}
		cm = cm.GetComByID(msg.CID)
		if cm.ID == 0 { //如果不是此adx的数据
			return true
		}
		cmp := index.CPINDEX.GetCompaign(msg.CID)

		if cmp.ID != 0 { //如果索引已存在则删除,保存历史纪录
			//保存历史纪录
			dailyBudgetRecores := cmp.DailyBudgetRecores
			dailyPPBRecords := cmp.DailyPPBRecords
			totalBudgetRecords := cmp.TotalBudgetRecords
			freqRecords := cmp.FreqRecords
			cm.DailyBudgetRecores = dailyBudgetRecores
			cm.DailyPPBRecords = dailyPPBRecords
			cm.TotalBudgetRecords = totalBudgetRecords
			cm.FreqRecords = freqRecords
			//删除老记录
			index.CPINDEX.Remove(msg.CID)
		}
		index.CPINDEX.Add(cm)

		index.CPINDEX.SetVersion(index.CPINDEX.GetVersion() + 1)

		// fmt.Println("更新index====", index.CPINDEX.GetCompaign(msg.CID).EndDate)

	case config.MQNDEL:
		index.CPINDEX.Remove(msg.CID)
		index.CPINDEX.SetVersion(index.CPINDEX.GetVersion() + 1)
	case config.MQNPENDING:
		index.CPINDEX.SetCompaignStatus(msg.CID, config.CompaignStatus[config.CS_PENDING])
	case config.MQNRUNNINT:
		index.CPINDEX.SetCompaignStatus(msg.CID, config.CompaignStatus[config.CS_RUNNING])
	}
	// fmt.Println(index.CPINDEX.GetCompaign(msg.CID).DailyBudget, index.CPINDEX.GetCompaign(msg.CID).DailyBudgetRecores.Cost)
	success = true
	return
}
