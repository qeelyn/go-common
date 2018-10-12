package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

const (
	// opentracing log key is trace.traceid
	ContextHeaderName = "qeelyn-traceid"
	TraceIdKey        = "traceid"
)

type Logger struct {
	zap        *zap.Logger
	sugar      *zap.SugaredLogger
	ToZapField func(values []interface{}) []zapcore.Field
}

func NewFileLogger(config map[string]interface{}) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	fp := config["filename"].(string)
	if _, err := os.Stat(fp); err != nil {
		if os.MkdirAll(filepath.Dir(fp), os.FileMode(0755)) != nil {
			panic(fmt.Errorf("Invalid Logger filename: %s", fp))
		}
		fi, err := os.Create(fp)
		defer fi.Close()
		if err != nil {
			panic(fmt.Errorf("create Logger filename: %s failure", fp))
		}
	}

	ms, _ := config["maxsize"].(int) // megabytes
	mbk, _ := config["maxbackups"].(int)
	ma, _ := config["maxsize"].(int) // days
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fp,
		MaxSize:    ms,
		MaxBackups: mbk,
		MaxAge:     ma,
	})

	level := zap.InfoLevel
	if cLevel, ok := config["level"].(int); ok {
		if cLevel >= -1 && cLevel <= 5 {
			level = levelConv(int8(cLevel))
		}
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		level,
	)
	return core
}

func NewStdLogger() zapcore.Core {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)
	return zapcore.NewCore(consoleEncoder, consoleDebugging, zap.DebugLevel)
}

func (l *Logger) SetZap(zap *zap.Logger) {
	l.zap = zap
	l.sugar = zap.Sugar()
}

// 类型化日志
func (l *Logger) Strict() *zap.Logger {
	return l.zap
}

// 语法糖方式,记录简单信息
func (l *Logger) Sugared() *zap.SugaredLogger {
	return l.sugar
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	return l.zap.With(TraceIdField(ctx))
}

// get trace id of zap field type
func TraceIdField(ctx context.Context) zap.Field {
	val, _ := ctx.Value(ContextHeaderName).(string)
	return zap.String(TraceIdKey, val)
}

func NewLogger(cores ...zapcore.Core) *Logger {
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core)
	return &Logger{zap: zapLogger, sugar: zapLogger.Sugar()}
}

func levelConv(level int8) zapcore.Level {
	switch level {
	case -1:
		return zap.DebugLevel
	case 0:
		return zap.InfoLevel
	case 1:
		return zap.WarnLevel
	case 2:
		return zap.ErrorLevel
	case 3:
		return zap.DPanicLevel
	case 4:
		return zap.PanicLevel
	case 5:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// Print passes arguments to Println
func (l *Logger) Print(values ...interface{}) {
	l.Println(values)
}

// for gorm & zap
func (l *Logger) Println(values []interface{}) {
	if l.ToZapField != nil {
		l.Strict().Info("gorm", l.ToZapField(values)...)
	}
}
