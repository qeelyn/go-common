package etcdv3

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/qeelyn/go-common/config/options"
	"github.com/spf13/viper"
	"io"
	"log"
	"time"
)

type etcdConfigProvider struct {
	Options  *options.Options
	client   *clientv3.Client
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
}

func Build(options *options.Options) error {
	provider, err := NewEtcdConfigProvider(options)
	if err != nil {
		return err
	}
	viper.RemoteConfig = provider
	return nil
}

func NewEtcdConfigProvider(options *options.Options) (*etcdConfigProvider, error) {
	cnf := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	if options.Timeout == 0 {
		options.Timeout = 5 * time.Second
	}

	if options.Secure || options.TLSConfig != nil {
		tlsConfig := options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		cnf.TLS = tlsConfig
	}

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		cnf.Endpoints = cAddrs
	}
	cnf.DialTimeout = options.Timeout
	cli, err := clientv3.New(cnf)
	if err != nil {
		return nil, err
	}
	e := &etcdConfigProvider{
		client:   cli,
		Options:  options,
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}

	return e, nil
}

func (t etcdConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	val, err := t.etcdGet(rp)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(val), nil
}

func (t etcdConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	val, err := t.etcdGet(rp)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(val), nil
}

func (t etcdConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	quit := make(chan bool)
	quitwc := make(chan bool)
	viperResponsCh := make(chan *viper.RemoteResponse)
	rch := t.client.Watch(context.Background(), rp.Path(), clientv3.WithPrefix())
	go func(rch *clientv3.WatchChan, vr chan<- *viper.RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				quit <- true
				return
			default:
				for n := range *rch {
					for _, ev := range n.Events {
						switch ev.Type {
						case mvccpb.PUT:
							viperResponsCh <- &viper.RemoteResponse{
								Error: n.Err(),
								Value: ev.Kv.Value,
							}
							log.Print(1)
						case mvccpb.DELETE:
							log.Print(2)
							quit <- true
						}
					}
				}
			}

		}

	}(&rch, viperResponsCh, quitwc, quit)
	return viperResponsCh, quit
}

func (t etcdConfigProvider) etcdGet(rp viper.RemoteProvider) ([]byte, error) {
	getResp, err := t.client.Get(context.Background(), rp.Path(), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if len(getResp.Kvs) == 0 {
		return nil, errors.New("key's value is empty")
	}
	return getResp.Kvs[0].Value, nil

}
