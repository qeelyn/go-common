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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qeelyn/go-common/auth"
	"github.com/qeelyn/go-common/grpcx/registry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	Name   string
	Option *serverOptions
}

func NewServer(name string, opts ...Option) *Server {
	srv := &Server{
		Name:   name,
		Option: &serverOptions{},
	}
	srv.Option.applyOption(opts...)
	return srv
}

// micro service stack support.
// if you use prometheus, the http server implement yourself
// example:
//   httpServer := &http.Server{
// 		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
// 		Addr: fmt.Sprintf("0.0.0.0:%d", 9092)
// 	 }
//
//	go func() {
//		if err := httpServer.ListenAndServe(); err != nil {
//			log.Fatal("Unable to start a http server.")
//		}
//	}()
func Micro(name string, opts ...Option) (*Server, error) {
	var err error
	sOptions := &serverOptions{}
	sOptions.applyOption(opts...)
	srv := &Server{
		Name:   name,
		Option: sOptions,
	}

	uins := []grpc.UnaryServerInterceptor{grpc_ctxtags.UnaryServerInterceptor()}
	sins := []grpc.StreamServerInterceptor{grpc_ctxtags.StreamServerInterceptor()}

	if sOptions.tracer != nil {
		uins = append(uins, grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(sOptions.tracer)))
		sins = append(sins, grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(sOptions.tracer)))
	}

	if sOptions.logger == nil {
		sOptions.logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		uins = append(uins, grpc_zap.UnaryServerInterceptor(sOptions.logger))
		sins = append(sins, grpc_zap.StreamServerInterceptor(sOptions.logger))
		grpc_zap.ReplaceGrpcLogger(sOptions.logger)
	}
	// tag and tracer must at first
	sOptions.unaryServerInterceptors = append(uins, sOptions.unaryServerInterceptors...)
	sOptions.streamServerInterceptors = append(sins, sOptions.streamServerInterceptors...)

	if sOptions.prometheus {
		sOptions.applyOption(WithUnaryServerInterceptor(grpc_prometheus.UnaryServerInterceptor))
		sOptions.applyOption(WithStreamServerInterceptor(grpc_prometheus.StreamServerInterceptor))

	}

	if sOptions.authFunc != nil {
		sOptions.applyOption(WithUnaryServerInterceptor(grpc_auth.UnaryServerInterceptor(sOptions.authFunc)))
		sOptions.applyOption(WithStreamServerInterceptor(grpc_auth.StreamServerInterceptor(sOptions.authFunc)))
	}

	sOptions.applyOption(WithUnaryServerInterceptor(grpc_recovery.UnaryServerInterceptor(
		grpc_recovery.WithRecoveryHandler(sOptions.Recovery))))
	sOptions.applyOption(WithStreamServerInterceptor(grpc_recovery.StreamServerInterceptor(
		grpc_recovery.WithRecoveryHandler(sOptions.Recovery))))

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

	if t.Option.prometheus {
		t.StartPrometheus(rpcSrv)
	}

	if t.Option.register != nil {
		node := &registry.Node{Id: t.Name, Address: t.Option.registryListen}
		if err = t.Option.register.Register(t.Option.registryServiceName, node); err != nil {
			return err
		}
		defer t.Option.register.Unregister(t.Option.registryServiceName, node)
	}

	log.Printf("%s tcp server will be ready for listening at:%s", t.Name, listen)
	return rpcSrv.Serve(lis)
}

func (t Server) StartPrometheus(rpcSrv *grpc.Server) {
	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(rpcSrv)
	// standalone http server
	if t.Option.prometheusListen != "" {
		// Register Prometheus metrics handler.
		httpServer := &http.Server{
			Handler: promhttp.Handler(),
			Addr:    t.Option.prometheusListen,
		}
		go func() {
			log.Printf("starting prometheus http server at:%s", httpServer.Addr)
			if err := httpServer.ListenAndServe(); err != nil {
				log.Fatal("Unable to start a http server.")
			}
		}()
	}
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
