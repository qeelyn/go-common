package etcdv3

import (
	"context"
	"log"
	"strings"
	"time"

	"crypto/tls"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/qeelyn/go-common/grpcx/registry"
	"google.golang.org/grpc/resolver"
	"path"
)

const PREFIX = "qeelyn"

type etcdv3Registry struct {
	client   *clientv3.Client
	cc       resolver.ClientConn
	options  registry.Options
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	config := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}

	var options registry.Options
	for _, o := range opts {
		o(&options)
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

		config.TLS = tlsConfig
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
		config.Endpoints = cAddrs
	}
	config.DialTimeout = options.Timeout
	cli, _ := clientv3.New(config)
	e := &etcdv3Registry{
		client:   cli,
		options:  options,
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}

	return e
}

func init() {
	registry.DefaultRegistry = NewRegistry
}

func (t *etcdv3Registry) Register(serviceName string, node *registry.Node, opts ...registry.RegisterOption) error {

	//var leaseNotFound bool
	//refreshing lease if existing
	leaseID, ok := t.leases[serviceName]
	if ok {
		if _, err := t.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}

			// lease not found do register
			//leaseNotFound = true
		}
	}

	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.options.Timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	var err error
	if options.TTL.Seconds() > 0 {
		lgr, err = t.client.Grant(ctx, int64(options.TTL.Seconds()))
		if err != nil {
			return err
		}
	}

	if lgr != nil {
		_, err = t.client.Put(ctx, nodePath(serviceName, node.Id), node.Address, clientv3.WithLease(lgr.ID))
	} else {
		_, err = t.client.Put(ctx, nodePath(serviceName, node.Id), node.Address)
	}
	if err != nil {
		return err
	}

	// save our leaseID of the service
	if lgr != nil {
		t.leases[serviceName] = lgr.ID
	}

	return nil
}

func nodePath(s, id string) string {
	service := strings.Replace(s, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(PREFIX, service, node)
}

func servicePath(s string) string {
	return path.Join(PREFIX, strings.Replace(s, "/", "-", -1))
}

func (t *etcdv3Registry) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {

	t.cc = cc
	key := servicePath("/" + target.Scheme + "/" + target.Endpoint + "/")
	go t.watch(key)

	return t, nil
}

// Close closes the resolver.
func (t etcdv3Registry) Close() {
	log.Println("Close")
}

func (t etcdv3Registry) Scheme() string {
	return PREFIX
}

func (t etcdv3Registry) ResolveNow(rn resolver.ResolveNowOption) {
	log.Println("ResolveNow") // TODO check
}

func (t *etcdv3Registry) watch(keyPrefix string) {
	var addrList []resolver.Address

	getResp, err := t.client.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	} else {
		for i := range getResp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix)})
		}
	}

	t.cc.NewAddress(addrList)

	rch := t.client.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					t.cc.NewAddress(addrList)
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					t.cc.NewAddress(addrList)
				}
			}
			//log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

// UnRegister remove service from etcd
func (t *etcdv3Registry) UnRegister(name string, addr string) {
	if t.client != nil {
		t.client.Delete(context.Background(), "/"+PREFIX+"/"+name+"/"+addr)
	}
}
