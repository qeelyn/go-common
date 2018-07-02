package registry

import (
	"crypto/tls"
	"time"
)

// The registry provides an interface for service discovery
type Registry interface {
	Register(serviceName string, node *Node, opts ...RegisterOption) error
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

var DefaultRegistry func(opts ...Option) Registry

type Option func(*Options)

type Options struct {
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	Addrs     []string
}

type RegisterOption func(*RegisterOptions)

type RegisterOptions struct {
	TTL time.Duration
}

func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Secure communication with the registry
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}
