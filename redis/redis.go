package redis

import (
	"github.com/go-redis/redis"
)

func NewRedisByMap(config map[string]interface{} ) *redis.Client {
	option := &redis.Options{}
	if addr,ok := config["addr"];ok {
		option.Addr = addr.(string)
	}
	if pwd,ok := config["password"];ok {
		option.Password = pwd.(string)
	}
	if db,ok := config["db"];ok {
		option.DB = db.(int)
	}
	if ps,ok := config["poolsize"];ok {
		option.PoolSize = ps.(int)
	}
	return redis.NewClient(option)
}
