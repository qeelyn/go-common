package memcache

import (
	"errors"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/qeelyn/go-common/cache"
	"reflect"
	"github.com/qeelyn/go-common/cache/internal"
	"github.com/qeelyn/go-common/cache/internal/util"
)

// Cache Memcache adapter.
type Cache struct {
	conn     *memcache.Client
	conninfo []string
	codec       cache.Codec
}

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

func NewMemcacheClient(config map[string]interface{}) (*memcache.Client,error) {
	if _, ok := config["addr"]; !ok {
		return nil,errors.New("config has no addr key")
	}
	conn := strings.Split(config["addr"].(string), ";")
	return memcache.New(conn...),nil
}

// Get get value from memcache.
func (t *Cache) Get(key string, dest interface{}) error {
	item, err := t.conn.Get(key)
	if err != nil {
		return err
	}

	kv := reflect.ValueOf(dest)
	tv := kv.Elem()
	switch tv.Kind() {
	case reflect.String:
		return inernal.Scan(item.Value, dest)
	case reflect.Float32, reflect.Float64:
		return inernal.Scan(item.Value, dest)
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inernal.Scan(item.Value, dest)
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Struct:
		return t.codec.Unmarshal(item.Value, dest)
	default:
		return t.codec.Unmarshal(item.Value, dest)
	}

	return nil
}

// GetMulti get value from memcache.
func (t *Cache) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	mv, err := t.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, v.Value)
		}
		return rv
	}
	for i := 0; i < size; i++ {
		rv = append(rv, err)
	}
	return rv
}

// Set put value to memcache.
func (t *Cache) Set(key string, val interface{}, timeout time.Duration) error {
	var err error
	item := memcache.Item{Key: key, Expiration: int32(timeout.Seconds())}

	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			item.Value = []byte(util.AsString(val))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			item.Value = []byte(util.AsString(val))
		case reflect.Float64:
			item.Value = []byte(util.AsString(val))
		case reflect.Float32:
			item.Value = []byte(util.AsString(val))
		case reflect.Bool:
			item.Value = []byte(util.AsString(val))
		default:
			if item.Value,err = t.codec.Marshal(val);err!=nil{
				return err
			}
		}
	}

	return t.conn.Set(&item)
}

// Delete delete value in memcache.
func (t *Cache) Delete(key string) error {
	return t.conn.Delete(key)
}

// Incr increase counter.
func (t *Cache) Incr(key string) error {
	_, err := t.conn.Increment(key, 1)
	return err
}

// Decr decrease counter.
func (t *Cache) Decr(key string) error {
	_, err := t.conn.Decrement(key, 1)
	return err
}

// IsExist check value exists in memcache.
func (t *Cache) IsExist(key string) bool {
	_, err := t.conn.Get(key)
	return !(err != nil)
}

// ClearAll clear all cached in memcache.
func (t *Cache) FlushAll() error {
	return t.conn.FlushAll()
}

// StartAndGC start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (t *Cache) StartAndGC(config map[string]interface{}) error {
	var err error
	if t.conn,err = NewMemcacheClient(config);err != nil {
		return err
	}
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}


