package kafka

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
	"jvole.com/influx/util"
)

type influxdb struct {
	addr      string //连接地址
	user      string //用户名
	passwd    string //密码
	database  string //数据库名
	precision string //精确度
	client    client.Client
	buff      uint16       //当前缓存数量
	data      []InfluxData //数据集合
	mtx       sync.RWMutex
}

var InfluxConn *influxdb

func NewInfluxdb(addr, user, password, db, precision string) *influxdb {
	// fmt.Printf("addr:%s,user:%s,pwd:%s\n", addr, user, password)
	// 创建 point batch
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: password,
	})

	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"name": "错误",
			"err":  err,
		}).Errorln("出错了")
	}
	indb := &influxdb{addr, user, password, db, precision, c, 0, make([]InfluxData, 0), *new(sync.RWMutex)}
	return indb
}

type InfluxData struct { //存放到influxdb中的数据结构
	Table  string
	Tags   map[string]string
	Fields map[string]interface{}
	Time   time.Time
}

var (
	Buffer  = uint16(1000) //缓存，达到一定数量后写数据库，可提高效率。默认值：1000
	MaxLine = 500          //查询时最多显示行数
)

/*
 写influx数据库
 @parm tags 标签相当于属性
 @parm fields 存储的字段集合，key value
 @parm table 表明
 @parm times 时间
 @return error
*/
func (idb *influxdb) Insert(tags map[string]string, fields map[string]interface{}, table string, times time.Time) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     idb.addr,
		Username: idb.user,
		Password: idb.passwd,
	})
	defer c.Close()
	if err != nil {
		return err
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  idb.database,
		Precision: idb.precision,
	})
	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"name": "错误",
			"err":  err,
		}).Errorln("创建bp出错了")
		return err
	}

	pt, err := client.NewPoint(table, tags, fields, times)
	if err != nil {
		util.Log.WithFields(logrus.Fields{
			"name": "错误",
			"err":  err,
		}).Errorln("添加point出错了")
		return err
	}
	bp.AddPoint(pt)

	//写数据库
	if err := c.Write(bp); err != nil { //这里会把所有有效数据添加
		util.Log.WithFields(logrus.Fields{
			"name":   "Inserts错误",
			"server": idb.addr,
			"err":    err,
		}).Errorln("写出错了")
		return err
	}
	return nil
}

/*
 关闭连接
 @return error
*/
func (idb *influxdb) Close() error {
	return idb.client.Close()
}
