package etcdv3_test

import (
	"fmt"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
	"net"
	"net/http"
	"testing"
	"time"
)

const testSvrName = "test"

func init() {
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(fmt.Sprint("failed to listen: %s", err))
	}
	defer lis.Close()
	go newMockServer(lis)
}

func newMockServer(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("a") != "1" {
			http.Error(w, "error", 500)
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	return http.Serve(l, mux)
}

func TestEtcdv3Registry_Build(t *testing.T) {
	rr := balancer.Get("round_robin")
	r := etcdv3.NewRegistry().(resolver.Builder)
	resolver.Register(r)

	_, err := grpc.Dial(testSvrName, grpc.WithBalancerName(rr.Name()), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

}

func TestEtcdv3Registry_Register(t *testing.T) {
	r := etcdv3.NewRegistry(registry.Timeout(5 * time.Minute))
	err := r.Register(testSvrName, &registry.Node{Id: testSvrName, Address: ":9999"})
	if err != nil {
		t.Error(err)
	}
}
