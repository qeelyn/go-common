package cache

import (
	"time"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v3"
)

type Cache interface {
	// get cached value by key.
	Get(key string,dest interface{}) error
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key and expire time.
	Set(key string, val interface{}, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	FlushAll() error
	// start gc routine based on config string settings.
	StartAndGC(config map[string]interface{}) error
}

// the remote cache server need serialize and derialize data
type CodecInterface interface {
	Marshal(interface{}) (interface{}, error)
	Unmarshal([]byte, interface{}) error
}

type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName string, config map[string]interface{}) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}

type Codec struct {
}

func (t *Codec) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (t *Codec) Unmarshal(b []byte, v interface{}) error {
	return msgpack.Unmarshal(b, v)
}


