package service

import (
	"encoding/json"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/k0kubun/pp"

	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/bidder"
	"jvole.com/dsp/config"
	"jvole.com/dsp/filter"
	"jvole.com/dsp/index"
	"jvole.com/dsp/model"
	"jvole.com/dsp/pool"
	"jvole.com/dsp/rabbitmq"
)

type exadsDSPBidder struct {
	logger *log.Logger
	lock   *sync.RWMutex
}

/*
		索引同步接口
	 @return []byte 索引数据
*/
func (s exadsDSPBidder) SyncIndex() []byte {
	return index.CPINDEX.GetIndexByte()
}

/*
		暂停某compaign投放 pending
	 @return bool true成功，false失败
*/
func (s exadsDSPBidder) StopBidByCID(cid uint32) bool {
	var res bool
	index.CPINDEX.SetCompaignStatus(cid, config.CompaignStatus[config.CS_PENDING])
	if index.CPINDEX.GetCompaign(cid).Status == config.CompaignStatus[config.CS_PENDING] {
		res = true
	}
	return res
}

/*
    adx 访问接口
   @return error
*/
func (s exadsDSPBidder) Bidder(body []byte, host string) string {
	defer func() {
		if err := recover(); err != nil {
			pp.Println("bidder意外崩溃", err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	offer := model.NewExadsADX(body, config.ADXCode[config.ADX], host)
	compaignres := &CompaignResult{make([]model.Compaign, 0)}

	//	分配任务
	cpbitmap := filter.LevelFilterComp(index.CPINDEX.GetIndexBitmap(), offer)
	num := int(cpbitmap.GetCardinality())

	sp := pool.NewSpool(num, num)
	sp.SetFinishCallback(s.finish)
	i := cpbitmap.Iterator()
	for i.HasNext() {
		var result uint32
		cid := i.Next()
		worksp := &pool.Worksp{
			Result: result,
			JobFunc: func() interface{} {
				return s.job(*&Jobs{offer, cid})
			},
			ResFunc: func(obj interface{}) error {
				return s.dealres(obj, compaignres)
			},
		}
		sp.AddTask(*worksp)
	}

	go func() { //超时中断任务
		time.Sleep(time.Millisecond * time.Duration(config.FilterTimeout))
		sp.Stop()
	}()

	sp.Start()

	// fmt.Println("length====", len(compaignres.CompaignArr))
	if len(compaignres.CompaignArr) == 0 { //如果没有
		return ""
	}

	//竞价算法
	cp := bidder.GetCompain(compaignres.CompaignArr)

	//生成返回结果
	result := offer.GetBidderResponse(cp)

	return result
}

/*
	adx回调接口，告知是否进价成功
	@return error
*/
func (s exadsDSPBidder) ADXNotify(notify rabbitmq.WinNotify) error {

	body, _ := json.Marshal(notify)
	//通知tracking
	rabbitmq.RabbitMQConn.Publish(config.TrackingEXName, config.TrackingRouteKey, body)
	//通知dsp
	rabbitmq.RabbitMQConn.Publish(config.DSPEXName, config.DSPRouteKey, body)
	return nil
}

func (s *exadsDSPBidder) finish() {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	level.Info(*s.logger).Log(
		"method", f.Name(),
		"file", path.Base(file),
		"line", line,
		"msg", "所有进程任务结束",
	)
}

func (s *exadsDSPBidder) job(job Jobs) interface{} {
	var cid uint32 = 0
	if filter.CPFiltercomp(index.CPINDEX.GetCompaign(job.CID), job.OF) {
		cid = job.CID
		// compaignArr = append(compaignArr, index.CPINDEX.Compaign[cid])
	}
	return cid
}

func (s *exadsDSPBidder) dealres(obj interface{}, compres *CompaignResult) (err error) {
	if cid, ok := obj.(uint32); ok {
		if cid > 0 {
			s.lock.Lock()
			compres.CompaignArr = append(compres.CompaignArr, index.CPINDEX.GetCompaign(cid))
			s.lock.Unlock()
		}
	}
	return
}

/*
返回服务实例
x
@return  struct
*/
func NewExadsServer(logger *log.Logger) *exadsDSPBidder {
	sb := &exadsDSPBidder{logger, new(sync.RWMutex)}

	return sb
}
