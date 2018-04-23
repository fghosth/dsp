package service

import (
	"github.com/go-kit/kit/log"

	"jvole.com/dsp/config"
	"jvole.com/dsp/model"
	"jvole.com/dsp/rabbitmq"
)

type DSPBidder interface {
	/*
	    adx 访问接口
	   @return error
	*/
	Bidder(data []byte, host string) string
	/*
		adx回调接口，告知是否进价成功
		@return error
	*/
	ADXNotify(notify rabbitmq.WinNotify) error
	/*
	    索引同步接口
	   @return []byte 索引数据
	*/
	SyncIndex() []byte
	/*
	    暂停某compaign投放 pending
	   @return bool true成功，false失败
	*/
	StopBidByCID(cid uint32) bool
}

type Jobs struct { //过滤任务
	OF  model.Offer
	CID uint32
}
type CompaignResult struct { //过滤后能参与竞价的compaign
	CompaignArr []model.Compaign
}

//根据code获得adx对象
//TODO 需手动增加adx服务
func GetService(code uint32, logger *log.Logger) (sev DSPBidder) {

	switch code {
	case config.ADXCode["smaato"]:
		sev = NewSmaatoServer(logger)
	case config.ADXCode["exads"]:
		sev = NewExadsServer(logger)
	case config.ADXCode["mgid"]:
		sev = NewMGIDServer(logger)
	case config.ADXCode["popcash"]:
		sev = NewPopcashServer(logger)
	case config.ADXCode["propellerads"]:
		sev = NewPropelleradsServer(logger)
	case config.ADXCode["adsterra"]:
		sev = NewAdsterraServer(logger)
	case config.ADXCode["zeropark"]:
		sev = NewZeroparkServer(logger)
	case config.ADXCode["popundertotal"]:
		sev = NewPopundertotalServer(logger)
	case config.ADXCode["hilltopads"]:
		sev = NewHilltopadsServer(logger)
	case config.ADXCode["eroadv"]:
		sev = NewEroadvServer(logger)
	}
	return
}

// func Works(result chan interface{}, job interface{}) {
//
// 	var cid uint32 = 0
// 	// time.Sleep(time.Second * 2)
// 	if st, ok := job.(Jobs); ok {
// 		if filter.CPFiltercomp(index.CPINDEX.Compaign[st.CID], st.OF) {
// 			cid = st.CID
// 		}
// 	}
// 	result <- cid
// }

// pool.WorkpoolConn.AddWorkers(func(job interface{}) {
// 	service.Works(Result, job)
// })
