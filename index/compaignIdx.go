package index

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/config"
	"jvole.com/dsp/model"
	"jvole.com/dsp/redis"
	"jvole.com/dsp/util"
)

var (
	redisdb = redis.RedisConn

	logger *log.Logger
)

type CompaignIdx struct {
	bitmap   *IndexMap
	compaign *sync.Map
	version  uint64 //版本号，正能增加
	lock     *sync.RWMutex
}

var CPINDEX = NewCompaignIdx()

func init() {
	logger = &util.KitLogger
	log.With(*logger, "component", "Index")
}

//启动时初始化索引
func (cmpIdx *CompaignIdx) SetupIndex() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		pp.Println("意外崩溃", err) // 这里的err其实就是panic传入的内容，55
	// 	}
	// }()
	//初始化adult索引数据
	for _, v := range config.AdultList {
		config.AdultBitmap.Add(util.Hashcode(v))
	}
	// 初始化索引，读取redis和本地，判断版本号。如果没有则从数据库建立索引

	if ok, _ := util.PathExists(config.CPIDXNAME); ok {
		CPINDEX.LoadDisk()
		level.Info(*logger).Log(
			"msg", "本地加载索引文件",
			"compaign数量", cmpIdx.GetCPLen(),
		)
		return
	}
	//从其他服务器加载index
	for _, url := range config.DSPSyncIndexServers {
		resu, err := util.DoBytesPost(url, []byte(`{"server":0}`))
		if err != nil {
			level.Info(*logger).Log(
				"msg", "从其他服务器加载index失败",
				"err", err,
			)
		}
		if string(resu) == "" { //服务器错误
			continue
		}
		errcode, _ := jsonparser.GetInt(resu, "errcode")
		if errcode != 0 { //获取数据错误
			continue
		}
		data := util.DoZlibUnCompress(resu)

		CPINDEX.SetIndexFromByte(data)
		if cmpIdx.GetCPLen() > 0 {
			level.Info(*logger).Log(
				"msg", "从服务器『"+url+"』加载索引成功",
				"compaign数量", cmpIdx.GetCPLen(),
			)
			return
		}
	}

	// var redisversion uint64
	// CPINDEX.LoadRedis()
	//
	// redisversion = CPINDEX.Version
	// if redisversion != 0 {
	// 	return
	// }
	data := []byte(`{"cmd":"` + config.DSPSQL["GetTotalCom"] + ` and (adExchanges & ` + strconv.Itoa(int(config.ADXCode[config.ADX])) + `)!=0 "}`)
	res, err := util.DoBytesPost(config.CompaignURL, data)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "获取总数量",
			"err", err,
		)
	}
	totalnum, _, _, _ := jsonparser.Get(res, "data", "[0]", "campaign", "total")

	total, _ := strconv.Atoi(string(totalnum)) //总记录数

	var page int //总页数
	page = total / config.PerNum
	if total%config.PerNum > 0 {
		page = page + 1
	}
	// fmt.Println(total, page)
	// return
	for i := 0; i < page; i++ { //分批查找数据建立索引
		data = []byte(`{"cmd":"` + config.DSPSQL["GetCompaign"] + ` and (adExchanges & ` + strconv.Itoa(int(config.ADXCode[config.ADX])) + `)!=0 limit ` + strconv.Itoa(config.PerNum) + ` offset ` + strconv.Itoa(i*config.PerNum) + `"}`)
		// pp.Println(string(data))
		comp, err := util.DoBytesPost(config.CompaignURL, data)
		if err != nil {
			pc, file, line, _ := runtime.Caller(1)
			f := runtime.FuncForPC(pc)
			level.Error(*logger).Log(
				"method", f.Name(),
				"file", path.Base(file),
				"line", line,
				"msg", "获取总数量",
				"err", err,
			)
		}

		c := &model.Compaign{}
		arrcomp := c.GetData(comp) //获得compaign
		for _, v := range arrcomp {
			CPINDEX.Add(v) //建立索引
		}
	}
	level.Info(*logger).Log(
		"msg", "从数据库加载index",
		"num of compaign", cmpIdx.GetCPLen(),
	)
}
func NewCompaignIdx() *CompaignIdx {
	cmpidx := &CompaignIdx{
		NewIndexMap(),
		new(sync.Map),
		0,
		new(sync.RWMutex),
		// rbt.NewWithIntComparator(),
		// NewCompainList(),
	}
	return cmpidx
}

//获得索引中所有compaign
func (cmpIdx *CompaignIdx) GetAllCompaigns() (cmp []model.Compaign) {
	cmpIdx.compaign.Range(func(key interface{}, value interface{}) bool {
		if v, ok := value.(model.Compaign); ok {
			cmp = append(cmp, v)
		}
		return true
	})
	return
}

//获得bitmap index
func (cmpIdx *CompaignIdx) GetIndexBitmap() (bm IndexMap) {
	cmpIdx.lock.RLock()
	defer cmpIdx.lock.RUnlock()
	bm = *cmpIdx.bitmap
	return
}

//获得compaign长度
func (cmpIdx *CompaignIdx) GetCPLen() (length int) {
	cmpIdx.lock.RLock()
	defer cmpIdx.lock.RUnlock()
	length = int(cmpIdx.bitmap.Compaign.GetCardinality())
	return
}

//改变compaign状态
func (cmpIdx *CompaignIdx) SetCompaignStatus(cid uint32, status uint8) {
	cmp, _ := cmpIdx.compaign.Load(cid)
	if v, ok := cmp.(model.Compaign); ok {
		v.Status = status
		if status == config.CompaignStatus[config.CS_PENDING] && (v.FilterSet&config.FilterCode["StatusFilter"] == 0) {
			v.FilterSet = v.FilterSet + config.FilterCode["StatusFilter"]
		}
		if status == config.CompaignStatus[config.CS_RUNNING] && (v.FilterSet&config.FilterCode["StatusFilter"] != 0) {
			v.FilterSet = v.FilterSet - config.FilterCode["StatusFilter"]
		}
		cmpIdx.compaign.Store(cid, v)
	}
}

//获取compaign状态
func (cmpIdx *CompaignIdx) GetCompaignStatus(cid uint32) (status uint8) {
	cmp, _ := cmpIdx.compaign.Load(cid)
	if v, ok := cmp.(model.Compaign); ok {
		status = v.Status
	}
	return
}

//设置版本号
func (cmpIdx *CompaignIdx) SetVersion(version uint64) {
	atomic.StoreUint64(&cmpIdx.version, version)
}

//获取版本号
func (cmpIdx *CompaignIdx) GetVersion() uint64 {
	return atomic.LoadUint64(&cmpIdx.version)
}

//添加索引记录
func (cmpIdx *CompaignIdx) Add(cp model.Compaign) {
	cmpIdx.compaign.Store(cp.ID, cp)
	cmpIdx.lock.Lock()
	defer cmpIdx.lock.Unlock()
	cmpIdx.bitmap.Add(cp)

	// cmpIdx.RBTree.Put(cp.Score, cp.ID)
	// cmpIdx.CompaignOrderByScore.AddCompain(cp)
}

//删除索引记录
func (cmpIdx *CompaignIdx) Remove(cid uint32) {
	cmpIdx.compaign.Delete(cid)
	cmpIdx.lock.Lock()
	defer cmpIdx.lock.Unlock()
	cmpIdx.bitmap.Remove(cid)

	// cmpIdx.RBTree.Remove(int(cid))
	// cmpIdx.CompaignOrderByScore.RemoveCompain(cid)
}

//获取compaign
func (cmpIdx *CompaignIdx) GetCompaign(cid uint32) (compaign model.Compaign) {
	cmp, _ := cmpIdx.compaign.Load(cid)
	if v, ok := cmp.(model.Compaign); ok {
		compaign = v
	}
	return compaign
}

//定时索引维护 更新记录
func (cmpIdx *CompaignIdx) IndexCheck() {
	go func() {
		for range time.Tick(time.Duration(config.Interval) * time.Second) {
			cmpIdx.compaign.Range(func(key interface{}, value interface{}) bool {
				if v, ok := value.(model.Compaign); ok {
					//遍历所有compaign
					rest, _ := util.CovnNOWUTC2Location(v.TimeZone.ZoneName)
					sinceTime, tillTime, _ := util.BgeinAndEndDAYOfZone(v.TimeZone.ZoneName)
					// pp.Println("indexcheck", int(v.ID), rest, v.EndDate)
					//=================日预算
					if !v.UnlimitBudget { //有预算限制
						//日预算记录
						if rest.After(v.EndDate) { //如果第二天了
							v.DailyBudgetRecores.Cost = 0
							v.DailyBudgetRecores.SinceTime = sinceTime
							v.DailyBudgetRecores.TillTime = tillTime
							//每日每广告位预算记录
							v.DailyPPBRecords.Offer = make(map[string]uint32, config.MapLength)
						}
					}
					//=================投放频率记录

					if v.FreqCapEnabled { //开启投放频率

						sinceTime = time.Now().UTC()
						interval := fmt.Sprint(v.FreqTimeWindow) + "h"
						d, _ := time.ParseDuration(interval)
						tillTime = sinceTime.Add(d)

						// pp.Println("indexcheck", int(v.ID), sinceTimef, v.FreqRecords.TillTime)
						if sinceTime.After(v.FreqRecords.TillTime) { //当前时间超过了之前的结束时间
							v.FreqRecords.SinceTime = sinceTime
							v.FreqRecords.TillTime = tillTime
							v.FreqRecords.User = make(map[string]uint32, config.MapLength)
							v.FreqRecords.Device = make(map[string]uint32, config.MapLength)
						}
					}

					cmpIdx.compaign.Store(key, v)

				}
				return true
			})

		}
	}()
}

//竞价成功后修改索引记录 cid CompaignIdx   t 竞价成功的UTC时间  price 价格
func (cmpIdx *CompaignIdx) BidSuccess(cid uint32, price uint32, position string, user, device string) (err error) {
	cmpIdx.lock.Lock()
	defer cmpIdx.lock.Unlock()
	var cm model.Compaign
	cmp, _ := cmpIdx.compaign.Load(cid)
	if v, ok := cmp.(model.Compaign); ok {
		cm = v
	} else {
		return errors.New("no compaign")
	}
	// fmt.Println(price, "===", cm.DailyBudgetRecores.Cost)
	//日预算记录
	if !cm.UnlimitBudget { //有预算限制
		cm.DailyBudgetRecores.Cost = cm.DailyBudgetRecores.Cost + price
		//总预算记录
		cm.TotalBudgetRecords.Cost = cm.TotalBudgetRecords.Cost + price
		//TODO 验证当花费高于总预算暂停投放 pending
		if cm.TotalBudgetRecords.Cost >= cm.TotalBudget { //当花费高于总预算暂停投放 pending
			cmpIdx.SetCompaignStatus(cid, config.CompaignStatus[config.CS_PENDING])
		}
	}
	// fmt.Println(price, "===", cm.DailyBudgetRecores.Cost)
	// pp.Println(position, "===", cm.DailyPPBRecords.Offer[position])
	if cm.DailyPerPlacementBudget > 0 { //有预算限制
		//每日每广告位预算记录
		_, ok := cm.DailyPPBRecords.Offer[position]
		if ok { //存在
			cm.DailyPPBRecords.Offer[position] = cm.DailyPPBRecords.Offer[position] + price
		} else { //不存在
			cm.DailyPPBRecords.Offer[position] = price
		}
	}
	// pp.Println(position, "===", cm.DailyPPBRecords.Offer[position])
	//投放频率记录
	// pp.Println("device", cm.FreqRecords.Device, "before===", "user", cm.FreqRecords.User)

	if cm.FreqCapEnabled && (device != "" || user != "") { //开启,并device和user不同时为空
		switch cm.FreqCapType {
		case 1: //设备
			_, ok := cm.FreqRecords.Device[device]
			if ok { //存在
				cm.FreqRecords.Device[device] = cm.FreqRecords.Device[device] + 1
			} else { //不存在
				cm.FreqRecords.Device[device] = 1
			}
		case 2: //用户
			_, ok := cm.FreqRecords.User[user]
			if ok { //存在
				cm.FreqRecords.User[user] = cm.FreqRecords.User[user] + 1
			} else { //不存在
				cm.FreqRecords.User[user] = 1
			}
		}
	}
	// pp.Println("device", cm.FreqRecords.Device, "after===", "user", cm.FreqRecords.User)
	cmpIdx.compaign.Store(cid, cm)
	//
	return
}

//清空索引
func (cmpIdx *CompaignIdx) Clear() {
	cmpIdx.compaign.Range(func(key interface{}, value interface{}) bool {
		cmpIdx.compaign.Delete(key)
		return true
	})
	cmpIdx.lock.Lock()
	defer cmpIdx.lock.Unlock()
	cmpIdx.bitmap.ClearAll()
	// cmpIdx.RBTree.Clear()
	// cmpIdx.CompaignOrderByScore.Clear()
}

//获得索引的byte
func (cmpIdx *CompaignIdx) GetIndexByte() []byte {
	data, _ := util.EncodeStructToByte(cmpIdx.GetAllCompaigns())
	return data
}

//从compaign数组恢复index
func (cmpIdx *CompaignIdx) RecoverIndex(data []model.Compaign) {
	for _, v := range data {
		// cmpIdx.compaign.Store(v.ID, v)
		cmpIdx.Add(v)
	}

}

//从byte加载到索引
func (cmpIdx *CompaignIdx) SetIndexFromByte(data []byte) error {
	var cmarr []model.Compaign
	err := util.DecodeByteToStruct(data, &cmarr)
	if err != nil {
		fmt.Println(err)
	}
	cmpIdx.RecoverIndex(cmarr)
	return err
}

//保存到redis
func (cmpIdx *CompaignIdx) SaveRedis() {
	data, _ := util.EncodeStructToByte(cmpIdx.GetAllCompaigns())
	data = util.DoZlibCompress(data)
	cmpIdx.lock.RLock()
	defer cmpIdx.lock.RUnlock()
	redisdb.Set(config.RedisCPVersionKey, util.Int64ToBytes(cmpIdx.version), 0) //设置当前版本号
	redisdb.Set(config.RedisIndexKey, data, 0)
}

//从redis读取
func (cmpIdx *CompaignIdx) LoadRedis() {
	cmpIdx.lock.RLock()
	data := redisdb.Get(config.RedisIndexKey)
	cmpIdx.lock.RUnlock()
	if data == nil {
		return
	}
	data = util.DoZlibUnCompress(data)
	var cmarr []model.Compaign
	err := util.DecodeByteToStruct(data, cmarr)
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
	cmpIdx.RecoverIndex(cmarr)
}

//保存磁盘
func (cmpIdx *CompaignIdx) SaveDisk() {

	data, err := util.EncodeStructToByte(cmpIdx.GetAllCompaigns())

	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "结构体转byte错误",
			"err", err,
		)
	}
	data = util.DoZlibCompress(data)
	cmpIdx.lock.RLock()
	err = ioutil.WriteFile(config.CPIDXNAME, data, 0666) //写入文件(字节数组)
	cmpIdx.lock.RUnlock()
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "保存磁盘错误",
			"err", err,
		)
	}

}

//读取磁盘缓存
func (cmpIdx *CompaignIdx) LoadDisk() {

	data := util.ReadFile(config.CPIDXNAME)
	if data == nil {
		return
	}
	data = util.DoZlibUnCompress(data)
	var cmarr []model.Compaign
	err := util.DecodeByteToStruct(data, &cmarr)
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
	// cmpIdx.lock.Unlock()
	cmpIdx.RecoverIndex(cmarr)
}
