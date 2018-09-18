package opentracing

import (
	"github.com/opentracing/opentracing-go"
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
