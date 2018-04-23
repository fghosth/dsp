package db

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"

	"github.com/buger/jsonparser"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/config"
	"jvole.com/dsp/util"
)

var logger log.Logger

func init() {
	logger = util.KitLogger
	log.With(logger, "component", "RabbitMQ")
}

type influxdb struct {
}

type insertRequest struct {
	Field struct {
		Data string `json:"data"`
		OID  string `json:"oid"`
	} `json:"field"`
	Table string `json:"table"`
	Tags  struct {
		CID    string `json:"cid"`
		Server string `json:"server"`
	} `json:"tags"`
	UID uint32 `json:"uid"`
}

type queryRequest struct {
	Cmd       string   `json:"cmd"`
	Db        string   `json:"db"`
	Limit     int      `json:"limit"`
	Offset    int      `json:"offset"`
	Precision string   `json:"precision"`
	UID       []uint32 `json:"uid"`
}

var InfluxdbConn = new(influxdb)

func (ifdb influxdb) Query(oid string, cid, uid uint32) []byte {
	qr := new(queryRequest)
	qr.Cmd = "select data from " + config.InfluxdbTable + " where server='" + config.Server + "' and cid='" + fmt.Sprint(cid) + "' and oid='" + oid + "'"
	qr.Limit = 1
	qr.Offset = 0
	qr.Precision = "ns"
	qr.UID = []uint32{uid}
	qrcmd, err := json.Marshal(qr)
	var data []byte
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "结构解析错误",
			"err", err,
		)
	}
	res, err := util.DoBytesPost(config.InfluxdbSelectURL, qrcmd)
	if err != nil { //调用接口失败
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "获取总数量",
			"err", err,
		)
	}
	Errcode, _ := jsonparser.GetInt(res, "Errcode")
	if Errcode != 0 { //查询失败
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "获取总数量",
			"err", err,
		)
		return data
	}
	//TODO 验证结构
	data, _, _, _ = jsonparser.Get(res, "Data", "[0]", "[0]", "values", "[0]", "[1]")
	return data
}

func (ifdb influxdb) Insert(oid, offer string, uid, cid uint32) {
	od := new(insertRequest)
	od.Field.Data = offer
	od.Table = config.InfluxdbDatabase
	od.Tags.CID = fmt.Sprint(cid)
	od.Tags.Server = config.Server
	od.UID = uid
	data, err := json.Marshal(od)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "结构解析错误",
			"err", err,
		)
	}
	res, err := util.DoBytesPost(config.InfluxdbInsertNowURL, data)
	if err != nil { //调用接口失败
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "执行插入",
			"err", err,
		)
	}
	succ, _ := jsonparser.GetInt(res, "Errcode")
	if succ != 0 { //插入失败
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "是否插入成功",
			"err", err,
		)
	}
}
