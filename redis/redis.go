package redis

import (
	"path"
	"runtime"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-redis/redis"
	"jvole.com/dsp/config"
	"jvole.com/dsp/util"
)

var logger *log.Logger

func init() {
	// fmt.Println(config.RedisURL, config.RedisPWD, config.RedisDB)
	logger = &util.KitLogger
	log.With(*logger, "component", "Redis")
}

type RedisDB struct {
	Host     string
	Password string
	DB       int
	Client   *redis.Client
}

var RedisConn = NewRedisDB(config.RedisURL, config.RedisPWD, config.RedisDB)

func NewRedisDB(host, pwd string, db int) *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})
	return &RedisDB{
		host,
		pwd,
		db,
		client,
	}
}

//获取
func (rdb RedisDB) Get(key string) []byte {
	val, err := rdb.Client.Get(key).Bytes()
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "获取rediskey失败",
			"err", err,
		)

	}
	return val
}

//保存
func (rdb RedisDB) Set(key string, value []byte, expiration time.Duration) error {
	err := rdb.Client.Set(key, value, expiration).Err()
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "设置rediskey失败",
			"err", err,
		)
	}
	return err
}

func (rdb RedisDB) Close() {
	rdb.Client.Close()
}
