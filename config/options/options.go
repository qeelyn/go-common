package options

import (
	"crypto/tls"
	"time"
)

type Option func(*Options)

type Options struct {
	// path for local or remote uri
	Path string
	// config fileName
	FileName string
	// remote host and port
	Addrs []string
	// enable secure
	Secure bool
	// tls config for remote
	TLSConfig *tls.Config
	Timeout   time.Duration
}
