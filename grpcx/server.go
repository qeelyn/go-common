package grpcx

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	Option *serverOptions
}

type serverOptions struct {
	tracer                  opentracing.Tracer
	logger                  *zap.Logger
	unaryServerInterceptors []grpc.UnaryServerInterceptor
	authFunc                grpc_auth.AuthFunc
}

func (t *serverOptions) applyOption(opts ...Option) *serverOptions {
	for _, v := range opts {
		v(t)
	}
	return t
}

type Option func(*serverOptions)

func WithLogger(logger *zap.Logger) Option {
	return func(options *serverOptions) {
		options.logger = logger
	}
}

func WithTracer(tracer opentracing.Tracer) Option {
	return func(options *serverOptions) {
		options.tracer = tracer
	}
}

func WithUnaryServerInterceptor(intercoptors ...grpc.UnaryServerInterceptor) Option {
	return func(options *serverOptions) {
		options.unaryServerInterceptors = append(options.unaryServerInterceptors, intercoptors...)
	}
}

func WithAuthFunc(authFunc grpc_auth.AuthFunc) Option {
	return func(options *serverOptions) {
		options.authFunc = authFunc
	}
}

func NewServer(opts ...Option) *Server {
	srv := &Server{
		Option: &serverOptions{},
	}
	srv.Option.applyOption(opts...)
	return srv
}

func Default(opts ...Option) (*Server, error) {
	var err error
	sOptions := &serverOptions{}

	sins := WithUnaryServerInterceptor(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
		grpc_prometheus.UnaryServerInterceptor)
	sOptions.applyOption(sins)

	sOptions.applyOption(opts...)
	if sOptions.authFunc != nil {
		sOptions.unaryServerInterceptors = append(sOptions.unaryServerInterceptors, grpc_auth.UnaryServerInterceptor(sOptions.authFunc))
	}
	// recovery at last
	sOptions.applyOption(WithUnaryServerInterceptor(grpc_recovery.UnaryServerInterceptor()))
	srv := &Server{
		Option: sOptions,
	}
	if srv.Option.logger == nil {
		srv.Option.logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}
	sOptions.unaryServerInterceptors = append(sOptions.unaryServerInterceptors, grpc_zap.UnaryServerInterceptor(srv.Option.logger))
	grpc_zap.ReplaceGrpcLogger(srv.Option.logger)
	return srv, err
}

func (t Server) BuildGrpcServer() *grpc.Server {
	var opts []grpc.ServerOption
	if len(t.Option.unaryServerInterceptors) > 0 {
		opts = append(opts, grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(t.Option.unaryServerInterceptors...),
		))
	}
	rpcSrv := grpc.NewServer(opts...)
	return rpcSrv
}

func (t Server) Run(rpcSrv *grpc.Server, listen string) error {
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return rpcSrv.Serve(lis)
}
