package grpcx_test

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/grpcx"
	"github.com/qeelyn/go-common/grpcx/dialer"
	"github.com/qeelyn/go-common/grpcx/internal/mock"
	"github.com/qeelyn/go-common/grpcx/internal/mock/prototest"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

func TestMicro(t *testing.T) {
	_, err := grpcx.Micro("test", grpcx.WithPrometheus(":9100"))
	if err != nil {
		t.Error(err)
	}
}

func TestMutilServer(t *testing.T) {
	a, err := grpcx.Micro("test", grpcx.WithPrometheus(":9100"))
	if err != nil {
		t.Error(err)
	}
	b, err := grpcx.Micro("test", grpcx.WithPrometheus(":9101"))
	if err != nil {
		t.Error(err)
	}
	arpc := a.BuildGrpcServer()
	//arpc.RegisterService(nil,nil)
	a.StartPrometheus(nil)
	go func() {
		a.Run(arpc, "9009")
	}()

	brpc := a.BuildGrpcServer()
	prototest.RegisterSayServer(brpc, &mock.Hello{})
	b.StartPrometheus(nil)
	b.Run(brpc, "9010")
}

func TestWithTracerLog(t *testing.T) {
	cfg := config.Configuration{}
	cfg.Headers = &jaeger.HeadersConfig{
		TraceContextHeaderName: grpc_opentracing.TagTraceId,
	}
	jLogger := jaeger.StdLogger
	jMetricsFactory := metrics.NullFactory
	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		config.Logger(jLogger),
		config.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	//cfg.Headers.TraceContextHeaderName = grpc_opentracing.TagTraceId
	defer closer.Close()

	a, err := grpcx.Micro("test",
		grpcx.WithTracer(opentracing.GlobalTracer()),
	)
	if err != nil {
		t.Error(err)
	}

	arpc := a.BuildGrpcServer()
	prototest.RegisterSayServer(arpc, &mock.Hello{})
	go a.Run(arpc, mock.TestSvrListen)

	cc, err := dialer.Dial(mock.TestSvrListen,
		dialer.WithTracer(opentracing.GlobalTracer()),
		dialer.WithDialOption(grpc.WithInsecure()),
	)
	if err != nil {
		panic(err)
	}
	opentracing.GlobalTracer().StartSpan("test")
	client := prototest.NewSayClient(cc)
	//ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	_, err = client.HelloPanic(context.Background(), &prototest.Request{})
	if err != nil {
		t.Error(err)
	}
}

func TestJwtAuthFunc(t *testing.T) {
	cnf := map[string]interface{}{
		"public-key":     nil,
		"encryption-key": "abc",
		"algorithm":      "HS256",
	}
	_, err := grpcx.Micro("test", grpcx.WithAuthFunc(grpcx.JwtAuthFunc(cnf)))
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithRegistry(t *testing.T) {
	var p1, p2 = ":80", "127.0.0.1:80"
	option := grpcx.WithRegistry(nil, "", p1)

	server := grpcx.NewServer("test", option)
	if server.Option.RegistryListen != p1 {
		t.Error()
	}

	os.Setenv("HOSTIP", "192.168.0.11")
	option = grpcx.WithRegistry(nil, "", p1)

	t.Log(os.Getenv("HOSTIP"))
	server = grpcx.NewServer("test", option)
	if server.Option.RegistryListen != os.Getenv("HOSTIP")+p1 {
		t.Error()
	}

	option = grpcx.WithRegistry(nil, "", p2)

	server = grpcx.NewServer("test", option)
	if server.Option.RegistryListen != p2 {
		t.Error()
	}
}
