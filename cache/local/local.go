package local

import (
	"reflect"
	gocache "github.com/patrickmn/go-cache"
	"time"
	"github.com/qeelyn/go-common/cache"
)

func NewLocalCache() cache.Cache {
	return &Cache{}
}

type Cache struct {
	localCache         *gocache.Cache
	localCacheDuration time.Duration
}

func (t *Cache) Get(key string, dest interface{}) error {
	if cacheData, ok := t.localCache.Get(key); ok {
		return convertAssign(dest, cacheData)
	}
	return nil
}

func (t *Cache) GetMulti(keys []string) []interface{} {
	var ret = []interface{}{}
	for _, v := range keys {
		if cacheData, ok := t.localCache.Get(v); ok {
			var val interface{}
			convertAssign(val, cacheData)
			ret = append(ret, val)
		} else {
			ret = append(ret, nil)
		}
	}
	return ret
}

func (t *Cache) Set(key string, val interface{}, expire time.Duration) error {
	t.localCache.Set(key, val, expire)
	return nil
}

func (t *Cache) Delete(key string) error {
	t.localCache.Delete(key)
	return nil
}

func (t *Cache) Incr(key string) error {
	if err := t.localCache.Increment(key, 1); err != nil {
		if t.localCache.Add(key, 1, 0) != nil {
			return t.localCache.Increment(key, 1)
		}
	}
	return nil
}

func (t *Cache) Decr(key string) error {
	if err := t.localCache.Decrement(key, 1); err != nil {
		if t.localCache.Add(key, -1, 0) != nil {
			return t.localCache.Decrement(key, 1)
		}
	}
	return nil
}

func (t *Cache) FlushAll() error {
	t.localCache.Flush()
	return nil
}

func (t *Cache) IsExist(key string) bool {
	if _, found := t.localCache.Items()[key]; found {
		return true
	}
	return false
}

func (t *Cache) StartAndGC(config map[string]interface{}) error {
	var defaultExp, cleanUp = 10*time.Minute, 30*time.Minute
	if duration, ok := config["duration"]; ok {
		defaultExp = time.Duration(duration.(int)) * time.Minute
	}
	if gc, ok := config["gc"]; ok {
		cleanUp = time.Duration(gc.(int)) * time.Minute
	}
	t.localCache = gocache.New(defaultExp, cleanUp)
	t.localCacheDuration = defaultExp
	return nil
}

func convertAssign(dest, src interface{}) error {
	dv := reflect.ValueOf(dest)
	sv := reflect.ValueOf(src)

	dk := dv.Elem()
	switch sv.Kind() {
	case reflect.Ptr:
		dk.Set(sv.Elem())
	default:
		dk.Set(sv)
	}
	return nil
}

func init() {
	cache.Register("local", NewLocalCache)
}
