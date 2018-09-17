package dialer

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//
type Options struct {
	Tracer                  opentracing.Tracer
	UnaryClientInterceptors []grpc.UnaryClientInterceptor
	DialOptions             []grpc.DialOption
}

type Option func(*Options)

// WithTracer traces rpc calls,the Tracer interceptor must put the last
// because the context  get value limit
// if use jeagertracer,the context header key will be uber-trace-id,it can't log by zap logger.
// if you want auto log by CtxTag,the key must be TagTraceID value,set the jeager's configuration key headers.TraceContextHeaderName to `trace.traceid`
func WithTracer(t opentracing.Tracer) Option {
	return func(options *Options) {
		options.Tracer = t
	}
}

func WithDialOption(gopts ...grpc.DialOption) Option {
	return func(options *Options) {
		options.DialOptions = gopts
	}
}

func WithAuth() grpc.UnaryClientInterceptor {
	return AuthUnaryInterceptor
}

// Dial returns a load balanced grpc client conn with tracing interceptor
func Dial(name string, opts ...Option) (*grpc.ClientConn, error) {
	options := Options{}

	for _, v := range opts {
		v(&options)
	}
	if options.Tracer != nil {
		// keep Tracer is last
		options.UnaryClientInterceptors = append(options.UnaryClientInterceptors, grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(options.Tracer)))
	}
	options.UnaryClientInterceptors = append(options.UnaryClientInterceptors, tracing.UnaryClientInterceptor(tracing.DefaultClientTraceIdFunc()))

	uopt := grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(options.UnaryClientInterceptors...))

	conn, err := grpc.Dial(name, append(options.DialOptions, uopt)...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", name, err)
	}

	return conn, nil
}

func WithUnaryClientInterceptor(interceptors ...grpc.UnaryClientInterceptor) Option {
	return func(options *Options) {
		options.UnaryClientInterceptors = append(options.UnaryClientInterceptors, interceptors...)
	}
}

func AuthUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	authHeader := ctx.Value("authorization")
	var md metadata.MD
	if authHeader != nil {
		md = metadata.Pairs("authorization", authHeader.(string))
	}
	orgHeader := ctx.Value("orgid")
	if orgHeader != nil {
		umd := metadata.Pairs("orgId", orgHeader.(string))
		md = metadata.Join(md, umd)
	}

	newCtx := metadata.NewOutgoingContext(ctx, md)

	return invoker(newCtx, method, req, reply, cc, opts...)
}
