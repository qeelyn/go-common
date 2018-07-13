package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	zap        *zap.Logger
	sugar      *zap.SugaredLogger
	ToZapField func(values []interface{}) []zapcore.Field
}

func (l *Logger) GetZap() *zap.Logger {
	return l.zap
}

func (l *Logger) SetZap(zap *zap.Logger) {
	l.zap = zap
	l.sugar = zap.Sugar()
}

func NewFileLogger(config map[string]interface{}) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	fp := config["filename"].(string)
	if _, err := os.Stat(fp); err != nil {
		if os.MkdirAll(filepath.Dir(fp), os.FileMode(0755)) != nil {
			panic(fmt.Errorf("Invalid Logger filename: %s", fp))
		}
		if _, err := os.Create(fp); err != nil {
			panic(fmt.Errorf("create Logger filename: %s failure", fp))
		}
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fp,
		MaxSize:    config["maxsize"].(int), // megabytes
		MaxBackups: config["maxbackups"].(int),
		MaxAge:     config["maxsize"].(int), // days
	})

	level := zap.InfoLevel
	if cLevel, ok := config["level"].(int); ok {
		if cLevel > -1 && cLevel < 5 {
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

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args)
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args)
}

// Print passes arguments to Println
func (l *Logger) Print(values ...interface{}) {
	l.Println(values)
}

// for gorm & zap
func (l *Logger) Println(values []interface{}) {
	if l.ToZapField != nil {
		l.GetZap().Info("gorm", l.ToZapField(values)...)
	}
}
