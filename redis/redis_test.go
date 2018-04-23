package redis

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"jvole.com/dsp/config"
)

func TestRedisDB_Set(t *testing.T) {
	type fields struct {
		Host     string
		Password string
		DB       int
		Client   *redis.Client
	}
	type args struct {
		key        string
		value      []byte
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"redisGet1", *&fields{config.RedisURL, config.RedisPWD, config.RedisDB, redis.NewClient(&redis.Options{Addr: config.RedisURL, Password: config.RedisPWD, DB: config.RedisDB})}, *&args{"key1", []byte("value1"), 0}, false},
		{"redisGet2", *&fields{config.RedisURL, config.RedisPWD, config.RedisDB, redis.NewClient(&redis.Options{Addr: config.RedisURL, Password: config.RedisPWD, DB: config.RedisDB})}, *&args{"key2", []byte("value2"), 0}, false},
		{"redisGet3", *&fields{config.RedisURL, config.RedisPWD, config.RedisDB, redis.NewClient(&redis.Options{Addr: config.RedisURL, Password: config.RedisPWD, DB: config.RedisDB})}, *&args{"key3", []byte("value3"), 0}, false},
	}
	for _, tt := range tests {
		rdb := RedisDB{
			Host:     tt.fields.Host,
			Password: tt.fields.Password,
			DB:       tt.fields.DB,
			Client:   tt.fields.Client,
		}
		if err := rdb.Set(tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
			t.Errorf("%q. RedisDB.Set() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestRedisDB_Get(t *testing.T) {
	type fields struct {
		Host     string
		Password string
		DB       int
		Client   *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{"redisGet1", *&fields{config.RedisURL, config.RedisPWD, config.RedisDB, redis.NewClient(&redis.Options{Addr: config.RedisURL, Password: config.RedisPWD, DB: config.RedisDB})}, *&args{"key1"}, []byte("value1")},
		{"redisGet2", *&fields{config.RedisURL, config.RedisPWD, config.RedisDB, redis.NewClient(&redis.Options{Addr: config.RedisURL, Password: config.RedisPWD, DB: config.RedisDB})}, *&args{"key2"}, []byte("value2")},
		{"redisGet3", *&fields{config.RedisURL, config.RedisPWD, config.RedisDB, redis.NewClient(&redis.Options{Addr: config.RedisURL, Password: config.RedisPWD, DB: config.RedisDB})}, *&args{"key3"}, []byte("value3")},
	}
	for _, tt := range tests {
		rdb := RedisDB{
			Host:     tt.fields.Host,
			Password: tt.fields.Password,
			DB:       tt.fields.DB,
			Client:   tt.fields.Client,
		}
		if got := rdb.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. RedisDB.Get() = %v, want %v", tt.name, string(got), string(tt.want))
		}
	}
}
