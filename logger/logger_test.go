package logger_test

import (
	"github.com/qeelyn/go-common/logger"
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	l := logger.NewLogger(logger.NewStdLogger())
	l.SetZap(l.Strict().With(zap.String("key", "keyValue")))
	l.Sugared().Error("test")
}
