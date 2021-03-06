package util

import (
	"errors"
	"time"

	"sync/cmd/app/options"

	"github.com/astaxie/beego/logs"
	"gopkg.in/redis.v5"
)

type RedisClient struct {
	*redis.Client
}

var (
	Client *RedisClient
)

func InitRedisClient() {
	cfg := options.GlobalConfig
	if cfg.Redis.RequiredPW {
		Client = &RedisClient{Client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
			PoolSize: cfg.PoolSize,
		})}
	}

	Client = &RedisClient{Client: redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})}

}

func (this *RedisClient) Set(key string, val interface{}) (string, error) {
	result, err := this.Client.Set(key, val, time.Minute*30).Result()
	if err != nil {
		logs.Error("set value err:%v", err)
		return "", err
	}
	return result, nil
}

func (this *RedisClient) Get(key string) (string, error) {
	result, err := this.Client.Get(key).Result()
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("key's value is not exsit")
	}
	return result, nil
}

func (this *RedisClient) MutilSet() {
	// this.Client.Pipelined(fn).
}
