package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"jvole.com/dsp/config"
	"jvole.com/dsp/dsperror"
	"jvole.com/dsp/index"
	"jvole.com/dsp/rabbitmq"
	"jvole.com/dsp/util"
)

type SyncIndexRequest struct {
	data []byte
}

func makeSyncIndexEndpoint(s DSPBidder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// req := request.(SyncIndexRequest)
		res := s.SyncIndex()
		// var err error
		var errcode int
		var msg string

		if len(res) == 0 {
			errcode = dsperror.ErrorSyncIndex
			msg = dsperror.ErrorText(errcode)
			// err = errors.New(msg)
			return &CommResponse{Errcode: errcode, Msg: msg, Data: nil}, nil
		}
		return res, nil
	}
}

type StopBidByCIDRequest struct {
	cid uint32
}

func makeStopBidByCIDEndpoint(s DSPBidder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StopBidByCIDRequest)

		res := s.StopBidByCID(req.cid)

		// var err error
		var errcode int
		var msg string
		if !res {
			errcode = dsperror.ErrorStopBid
			msg = dsperror.ErrorText(errcode)
			// err = errors.New(msg)
		}

		return &CommResponse{Errcode: errcode, Msg: msg, Data: index.CPINDEX.GetCompaign(req.cid).TotalBudgetRecords}, nil
		// return &CommResponse{Errcode: errcode, Msg: msg, Data: res}, nil

	}
}

type BidderRequest struct {
	data []byte
	host string
}

type NilError struct{ errMsg string }

func makeBidderEndpoint(s DSPBidder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(BidderRequest)
		if config.ADXReqTrackISON { //adxrequest日志是否开启
			util.Tracefile(string(req.data), config.ADXReqTrackPath)
		}
		bidresp := s.Bidder(req.data, req.host)
		if bidresp == "" { //不竞价返回204
			return 204, errNotBid
		}
		return bidresp, nil

	}
}

type ADXNotifyRequest struct {
	ID      string  //广告位id
	Price   float64 //竞价价格
	CID     uint32  //compaignid
	UID     uint32  //用户id
	Postion string  //广告位置标识
	Device  string  //设备标识
	T       int64   //竞价成功时间
	User    string  //广告端user标识
	ClickID string
}

func makeADXNotifyEndpoint(s DSPBidder) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ADXNotifyRequest)
		notify := &rabbitmq.WinNotify{
			ID:      req.ID,
			Price:   uint32(req.Price / float64(config.BidPriceFix) * float64(config.Cashfix)),
			CID:     req.CID,
			UID:     req.UID,
			Postion: req.Postion,
			Device:  req.Device,
			T:       req.T,
			User:    req.User,
			ADX:     config.ADXCode[config.ADX],
			ClickID: req.ClickID,
		}
		if config.ADXReqTrackWin { //adxrequest日志是否开启
			util.Tracefile(fmt.Sprint(notify), config.ADXReqTrackWinPath)
		}
		err := s.ADXNotify(*notify)
		if err != nil {
			return BIDResponse{err.Error()}, nil
		}
		return "ok", nil
	}
}

type CommResponse struct {
	Errcode int         `json:"errcode"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type BIDResponse struct {
	Result string
}
