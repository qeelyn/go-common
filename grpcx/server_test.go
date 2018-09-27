package grpcx_test

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/grpcx"
	"github.com/qeelyn/go-common/grpcx/authfg"
	"github.com/qeelyn/go-common/grpcx/dialer"
	"github.com/qeelyn/go-common/grpcx/internal/mock"
	"github.com/qeelyn/go-common/grpcx/internal/mock/prototest"
	"github.com/qeelyn/go-common/grpcx/tracing"
	logger2 "github.com/qeelyn/go-common/logger"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"os"
	"testing"
	"time"
)

func TestMicro(t *testing.T) {
	_, err := grpcx.Micro("test", grpcx.WithPrometheus(":9100"),
		grpcx.WithGrpcOption(grpc.KeepaliveParams(keepalive.ServerParameters{Time: 10 * time.Minute})),
	)
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
	if err == nil {
		t.Error(err)
	}
}

func TestJwtAuthFunc(t *testing.T) {
	cnf := map[string]interface{}{
		"public-key":     nil,
		"encryption-key": "abc",
		"algorithm":      "HS256",
	}
	_, err := grpcx.Micro("test", grpcx.WithAuthFunc(authfg.ServerJwtAuthFunc(cnf)))
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

func TestWithTracerId(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}

	a, err := grpcx.Micro("test",
		grpcx.WithLogger(logger),
		grpcx.WithUnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			v := metautils.ExtractIncoming(ctx).Get(logger2.ContextHeaderName)
			if v == "" {
				return nil, errors.New("server not receive trace context")
			}
			return handler(ctx, req)
		}),
	)
	if err != nil {
		t.Error(err)
	}

	arpc := a.BuildGrpcServer()
	prototest.RegisterSayServer(arpc, &mock.Hello{})
	go a.Run(arpc, mock.TestSvrListen)

	cc, err := dialer.Dial(mock.TestSvrListen,
		dialer.WithDialOption(grpc.WithInsecure()),
		dialer.WithTraceIdFunc(tracing.DefaultClientTraceIdFunc(true)),
	)
	if err != nil {
		panic(err)
	}

	client := prototest.NewSayClient(cc)
	//ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()
	ctx = tracing.ToContext(ctx, "testtraceid")
	//ctx = metadata.AppendToOutgoingContext(ctx,"trace.traceid","testtraceid")
	_, err = client.Hello(ctx, &prototest.Request{})
	if err != nil {
		t.Error(err)
	}
	// not set value,because interceptor return err
	ctx = context.Background()
	_, err = client.Hello(ctx, &prototest.Request{})
	if err == nil {
		t.Error("service error not found")
	}
}
