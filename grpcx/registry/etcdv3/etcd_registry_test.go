package etcdv3_test

import (
	"context"
	"fmt"
	"github.com/qeelyn/go-common/grpcx/internal/mock"
	"github.com/qeelyn/go-common/grpcx/internal/mock/prototest"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"testing"
	"time"
)

func init() {
	go mock.NewMockDiscoveryServer(":12345", "default")
}

func TestEtcdv3Registry_Build(t *testing.T) {
	r, _ := etcdv3.NewRegistry()
	b := r.(resolver.Builder)
	resolver.Register(b)

	_, err := grpc.Dial(mock.TestSvrName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

}

func TestEtcdv3Registry_Register(t *testing.T) {
	r, err := etcdv3.NewRegistry(registry.Timeout(5 * time.Minute))
	if err != nil {
		t.Error(err)
	}
	err = r.Register(mock.TestSvrName, &registry.Node{Id: mock.TestSvrName, Address: ":12345"})
	if err != nil {
		t.Error(err)
	}
}

func TestMulti(t *testing.T) {
	//
	go mock.NewMockDiscoveryServer(":12346", "n2")
	go mock.NewMockDiscoveryServer(":12347", "n3")
	r, _ := etcdv3.NewRegistry(registry.Timeout(30 * time.Second))
	b := r.(resolver.Builder)
	resolver.Register(b)
	//d1
	d1, err := grpc.Dial(b.Scheme()+"://author/"+mock.TestSvrName, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	client := prototest.NewSayClient(d1)
	if err != nil {
		panic(err)
	}
	for {
		out, err := client.Hello(context.Background(), &prototest.Request{Name: "test"})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(out)
		}
		<-time.After(time.Second)
	}
}
