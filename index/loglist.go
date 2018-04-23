package index

import (
	"container/list"
	"path"
	"runtime"
	"sync"

	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/util"
)

//索引变更日志
type LogList struct {
	List    *list.List
	Version uint64
	lock    *sync.RWMutex
}

type IndexLog struct {
	Cid     uint32 //compaignID
	Action  uint8  //增加compaign：1，修改compaign：2，删除compaign：3
	Version uint64
}

var (
	redisLogListKey    = "DSP:loglist"
	redisLogVersionKey = "DSP:logversion"
)

//增加log
func (loglist *LogList) Add(il IndexLog) {
	loglist.lock.Lock()
	loglist.List.PushBack(il)
	loglist.Version = il.Version
	loglist.lock.Unlock()
}

//获取log长度
func (loglist *LogList) Len() int {
	loglist.lock.RLock()
	defer loglist.lock.RUnlock()
	return loglist.List.Len()

}

//清空loglist
func (loglist *LogList) Clear() {
	loglist.lock.Lock()
	loglist.List.Init()
	loglist.Version = 0
	loglist.lock.Unlock()
}

//保存到redis
func (loglist *LogList) SaveToRedis() {

	loglist.lock.Lock()
	data, _ := util.EncodeStructToByte(loglist)
	data = util.DoZlibCompress(data)
	redisdb.Set(redisLogListKey, data, 0)
	redisdb.Set(redisLogVersionKey, util.Int64ToBytes(loglist.Version), 0) //设置当前版本号
	loglist.lock.Unlock()
}

//从redis中读取
func (loglist *LogList) LoadFromRedis() {
	loglist.lock.Lock()
	data := redisdb.Get(redisLogListKey)
	data = util.DoZlibUnCompress(data)
	err := util.DecodeByteToStruct(data, loglist)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "byte转结构体错误",
			"err", err,
		)
	}
	loglist.lock.Unlock()
}
