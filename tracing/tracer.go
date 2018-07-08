package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/gin-contrib/ginzap"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"log"
)

func NewTracer(config *jaegercfg.Configuration, serviceName string) opentracing.Tracer {
	if serviceName != "" {
		config.ServiceName = serviceName
	}
	jMetricsFactory := metrics.NullFactory
	tracer, _, err := config.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Panicf("cannot initialize Jaeger Tracer", zap.Error(err))
	}
	return tracer
}

type jaegerLoggerAdapter struct {
	logger *ginzap.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
