package loginflux

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

//日志插件，可以把日志记录到influxdb，实现了io.writer接口能够用于所有，传入io.writer的log日志插件（只支持json格式的）
type Loginflux struct {
	data   []logdata //数据集合
	server InfluxServer
	mtx    sync.RWMutex
}

type InfluxServer struct {
	Addr      string            //连接地址
	User      string            //用户名
	Passwd    string            //密码
	Database  string            //数据库名
	Table     string            //表名
	Precision string            //精确度
	Buff      uint16            //缓存数量
	Tags      map[string]string //存放influxdb时的tags
	KeyS      []string          //如果日志中出现这些字段就把它假如到tags中
}

type logdata struct { //存放到influxdb中的数据结构
	table  string
	tags   map[string]string
	fields map[string]interface{}
	time   time.Time
}

var (
	PrintScreen = false //是同时否输出到屏幕 true 为显示
)

//实现io.Writer
func (lf *Loginflux) Write(p []byte) (n int, err error) {
	lf.mtx.Lock()
	defer lf.mtx.Unlock()
	var result map[string]interface{}
	if err = json.Unmarshal(p, &result); err != nil {
		return n, err
	}
	for _, v := range lf.server.KeyS { //添加tags
		if val, ok := result[v].(string); ok {
			lf.server.Tags[v] = val
			delete(result, v) //去除tags中出现fields中也出现的重复内容
		}
	}
	ldata := &logdata{
		lf.server.Table,
		lf.server.Tags,
		result,
		time.Now(),
	}
	lf.data = append(lf.data, *ldata)
	if uint16(len(lf.data)) >= lf.server.Buff { //缓存满了写数据库
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     lf.server.Addr,
			Username: lf.server.User,
			Password: lf.server.Passwd,
		})
		defer c.Close() //关闭连接
		if err != nil {
			fmt.Println("创建客户端错误:" + err.Error())
		}
		// 创建 point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  lf.server.Database,
			Precision: lf.server.Precision,
		})
		if err != nil {
			fmt.Println("创建bp错误:" + err.Error())
		}
		for _, v := range lf.data {
			pt, err := client.NewPoint(v.table, v.tags, v.fields, v.time)
			if err != nil {
				fmt.Println("创建pt错误:" + err.Error())
			}
			bp.AddPoint(pt)
		}
		//写数据库
		if err := c.Write(bp); err != nil { //这里会把所有有效数据添加
			fmt.Println("日志写入错误:" + err.Error())
		} else {
			lf.data = nil //清除缓存
		}
	}

	//输出到屏幕
	if PrintScreen {
		os.Stdout.Write(p)
	}
	return n, err
}

func NewLoginflux(is InfluxServer) io.Writer {
	lf := &Loginflux{}
	lf.server = is
	return lf
}
