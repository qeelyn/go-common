package logger_test

import (
	"context"
	"github.com/qeelyn/go-common/logger"
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	l := logger.NewLogger(logger.NewStdLogger())
	l.SetZap(l.Strict().With(zap.String("key", "keyValue")))
	l.Sugared().Error("test")
}

func TestTraceIdField(t *testing.T) {
	ctx := context.Background()
	f := logger.TraceIdField(ctx)
	if f.String != "" {
		t.Fatal()
	}
}
