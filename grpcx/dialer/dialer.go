package dialer

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithTracer traces rpc calls
func WithTracer(t opentracing.Tracer) grpc.UnaryClientInterceptor {
	return otgrpc.OpenTracingClientInterceptor(t)
}

func WithAuth() grpc.UnaryClientInterceptor {
	return AuthUnaryInterceptor
}

// Dial returns a load balanced grpc client conn with tracing interceptor
func Dial(name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(name, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", name, err)
	}

	return conn, nil
}

func WithUnaryClientInterceptor(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...))
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

func newUserContext(ctx context.Context) context.Context {
	var md metadata.MD
	if uid := ctx.Value("userid"); uid != nil {
		md = metadata.Pairs("userid", uid.(string))
	}
	if oid := ctx.Value("orgid"); oid != nil {
		umd := metadata.Pairs("orgId", oid.(string))
		md = metadata.Join(md, umd)
	}

	rpcCtx := metadata.NewOutgoingContext(ctx, md)
	return rpcCtx
}
