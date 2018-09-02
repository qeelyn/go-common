package options

import "github.com/qeelyn/go-common/grpcx/registry"

type Option func(*Options)

type Options struct {
	// path for local or remote uri
	Path string
	// config fileName
	FileName string
	// registry center
	Registry registry.Registry
	// address of registry
	Addr string
}
