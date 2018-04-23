package cache

import (
	"path"
	"runtime"
	"sync"

	"github.com/coocood/freecache"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/config"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

func init() {
	logger = util.KitLogger
	log.With(logger, "component", "Cache")
}

//投放广告的缓存
type offerCache struct {
	OCache *freecache.Cache
	lock   *sync.RWMutex
}

var (
	OfferCacheConn = NewOfferCache()
	logger         log.Logger
)

func NewOfferCache() *offerCache {
	c := freecache.NewCache(config.CacheSize)
	return &offerCache{c, new(sync.RWMutex)}
}

func (oc *offerCache) Set(key []byte, offer model.Offer, expire int) {
	data := util.DoZlibCompress(offer.GetData())
	oc.lock.Lock()
	defer oc.lock.Unlock()
	oc.OCache.Set(key, data, expire)
}

func (oc *offerCache) Remove(key []byte) {
	oc.lock.Lock()
	defer oc.lock.Unlock()
	oc.OCache.Del(key)
}

//得到广告
func (oc *offerCache) Get(key []byte) []byte {
	oc.lock.RLock()
	defer oc.lock.RUnlock()
	var data []byte
	got, err := oc.OCache.Get(key)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "获取缓存失败",
			"err", err,
		)
	} else {
		data = util.DoZlibUnCompress(got)
	}
	return data
}
