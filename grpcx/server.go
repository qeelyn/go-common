package grpcx

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
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
		grpc_prometheus.UnaryServerInterceptor)
	sOptions.applyOption(sins)

	sOptions.applyOption(opts...)

	if sOptions.tracer != nil {
		usi := grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(sOptions.tracer))
		sOptions.unaryServerInterceptors = append([]grpc.UnaryServerInterceptor{usi}, sOptions.unaryServerInterceptors...)
	}

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

func AuthFunc(keyFile string) grpc_auth.AuthFunc {
	validator := auth.BearTokenValidator{
		PubKeyFile: keyFile,
		IdentityHandler: func(ctx context.Context, claims jwt.MapClaims) (*auth.Identity, error) {
			orgIdStr := metautils.ExtractIncoming(ctx).Get("orgid")
			id, _ := strconv.Atoi(claims["sub"].(string))
			orgId, _ := strconv.Atoi(orgIdStr)
			identity := &auth.Identity{
				Id:    int32(id),
				OrgId: int32(orgId),
			}
			return identity, nil
		},
	}
	validator.Init()
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return ctx, err
		}

		if id, err := validator.Validate(ctx, token); err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		} else {
			return context.WithValue(ctx, "user", id), nil
		}
	}
}
