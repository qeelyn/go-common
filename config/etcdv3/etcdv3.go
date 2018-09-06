package etcdv3

import (
	"bytes"
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/qeelyn/go-common/config/options"
	"github.com/spf13/viper"
	"io"
	"log"
)

type etcdConfigProvider struct {
	Options *options.Options
	client  *clientv3.Client
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
	if options.Registry == nil {
		return nil, errors.New("registry is not set")
	}
	client, ok := options.Registry.GetClient().(*clientv3.Client)
	if !ok {
		return nil, errors.New("registry client is not an etcd v3 client")
	}
	options.Addr = client.Endpoints()[0]
	e := &etcdConfigProvider{
		client:  client,
		Options: options,
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
