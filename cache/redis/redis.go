package redis

import (
	"github.com/go-redis/redis"
	"github.com/qeelyn/go-common/cache"
	qredis "github.com/qeelyn/go-common/redis"
	"time"
	"reflect"
	"github.com/qeelyn/go-common/cache/internal"
)

type Cache struct {
	redisClient *redis.Client
	codec       *cache.Codec
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{}
}

func (t *Cache) Get(key string, dest interface{}) error {
	var err error

	var data []byte
	if data, err = t.redisClient.Get(key).Bytes(); err != nil {
		return err
	}

	kv := reflect.ValueOf(dest)
	tv := kv.Elem()
	switch tv.Kind() {
	case reflect.String:
		return inernal.Scan(data, dest)
	case reflect.Float32, reflect.Float64:
		return inernal.Scan(data, dest)
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inernal.Scan(data, dest)
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Struct:
		return t.codec.Unmarshal(data, dest)
	default:
		return t.codec.Unmarshal(data, dest)
	}
}

// the int,string return origin,other need decode
func (t *Cache) GetMulti(keys []string) []interface{} {
	var args []interface{}
	for _, key := range keys {
		args = append(args, key)
	}
	values, err := t.redisClient.MGet(keys...).Result()
	if err != nil {
		return nil
	}
	return values
}

func (t *Cache) Put(key string, val interface{}, expire time.Duration) error {
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.String:
		return t.redisClient.Set(key, val, expire).Err()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return t.redisClient.Set(key, val, expire).Err()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return t.redisClient.Set(key, val, expire).Err()
	case reflect.Float64:
		return t.redisClient.Set(key, val, expire).Err()
	case reflect.Float32:
		return t.redisClient.Set(key, val, expire).Err()
	case reflect.Bool:
		return t.redisClient.Set(key, val, expire).Err()
	default:
		if data, err := t.codec.Marshal(val); err != nil {
			return err
		} else {
			return t.redisClient.Set(key, data, expire).Err()
		}
	}
}

func (t *Cache) Delete(key string) error {
	return t.redisClient.Del(key).Err()
}

func (t *Cache) Incr(key string) error {
	return t.redisClient.Incr(key).Err()
}

func (t *Cache) Decr(key string) error {
	return t.redisClient.Decr(key).Err()
}

func (t *Cache) FlushAll() error {
	return t.redisClient.FlushAll().Err()
}

func (t *Cache) IsExist(key string) bool {
	return t.redisClient.Exists(key).Val() != 0
}

func (t *Cache) StartAndGC(config map[string]interface{}) error {
	t.redisClient = qredis.NewRedisByMap(config)
	t.codec = &cache.Codec{}
	return nil
}

func init() {
	cache.Register("redis", NewRedisCache)
}
