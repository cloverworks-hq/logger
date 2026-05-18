package logger

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Trace(msg ...interface{})
	Tracef(format string, args ...interface{})
	Debug(msg ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg ...interface{})
	Infof(format string, args ...interface{})
	Warn(msg ...interface{})
	Warnf(format string, args ...interface{})
	Error(msg ...interface{})
	Errorf(format string, args ...interface{})
	ErrorWithErr(msg string, err error, fields ...map[string]interface{})
	Fatal(msg ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(msg ...interface{})
	Panicf(format string, args ...interface{})
	
	// With returns a new Logger with the provided fields attached.
	With(fields map[string]interface{}) Logger
}

type zapLogger struct {
	log *zap.Logger
}

type contextKey string

const loggerContextKey contextKey = "logger"

// New creates a new Logger configured with the provided Options.
func New(opts ...Option) Logger {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	var level zapcore.Level
	switch strings.ToLower(options.Level) {
	case "trace", "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	output := options.Output
	if output == "" {
		output = "stdout"
	}

	logConfig := zap.Config{
		OutputPaths: []string{output},
		Level:       zap.NewAtomicLevelAt(level),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:      "level",
			TimeKey:       "time",
			MessageKey:    "msg",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
			StacktraceKey: "stacktrace",
		},
		DisableStacktrace: false,
	}

	zl, err := logConfig.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{log: zl}
}

// FromContext retrieves a Logger from the context. If none is found, it returns a no-op logger or default logger.
// In a robust application, you typically initialize a base logger and attach it to the root context.
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return New() // Fallback
	}
	if l, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return l
	}
	return New() // Fallback
}

// WithContext returns a new context with the provided Logger attached.
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

func (l *zapLogger) Trace(msg ...interface{}) {
	l.log.Debug(fmt.Sprint(msg...))
}

func (l *zapLogger) Tracef(format string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Debug(format)
	} else {
		l.log.Debug(fmt.Sprintf(format, args...))
	}
}

func (l *zapLogger) Debug(msg ...interface{}) {
	l.log.Debug(fmt.Sprint(msg...))
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Debug(format)
	} else {
		l.log.Debug(fmt.Sprintf(format, args...))
	}
}

func (l *zapLogger) Info(msg ...interface{}) {
	l.log.Info(fmt.Sprint(msg...))
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Info(format)
	} else {
		l.log.Info(fmt.Sprintf(format, args...))
	}
}

func (l *zapLogger) Warn(msg ...interface{}) {
	l.log.Warn(fmt.Sprint(msg...))
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Warn(format)
	} else {
		l.log.Warn(fmt.Sprintf(format, args...))
	}
}

func (l *zapLogger) Error(msg ...interface{}) {
	l.log.Error(fmt.Sprint(msg...))
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Error(format)
	} else {
		l.log.Error(fmt.Sprintf(format, args...))
	}
}

// ErrorWithErr provides a structured way to log a message alongside an actual error object and optional fields.
func (l *zapLogger) ErrorWithErr(msg string, err error, fields ...map[string]interface{}) {
	var zapFields []zap.Field
	zapFields = append(zapFields, zap.NamedError("error", err))

	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}

	l.log.Error(msg, zapFields...)
	l.log.Sync()
}

func (l *zapLogger) Fatal(msg ...interface{}) {
	l.log.Fatal(fmt.Sprint(msg...))
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	if len(args) == 0 {
		l.log.Fatal(format)
	} else {
		l.log.Fatal(fmt.Sprintf(format, args...))
	}
}

func (l *zapLogger) Panic(msg ...interface{}) {
	msgStr := fmt.Sprint(msg...)
	l.log.Info(msgStr)
	panic(msgStr)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	var msgStr string
	if len(args) == 0 {
		msgStr = format
	} else {
		msgStr = fmt.Sprintf(format, args...)
	}
	l.log.Info(msgStr)
	panic(msgStr)
}

func (l *zapLogger) With(fields map[string]interface{}) Logger {
	zapFields := make([]zap.Field, 0, len(fields))

	for key, value := range fields {
		switch v := value.(type) {
		case string:
			zapFields = append(zapFields, zap.String(key, v))
		case int:
			zapFields = append(zapFields, zap.Int(key, v))
		case int64:
			zapFields = append(zapFields, zap.Int64(key, v))
		case float64:
			zapFields = append(zapFields, zap.Float64(key, v))
		case bool:
			zapFields = append(zapFields, zap.Bool(key, v))
		default:
			zapFields = append(zapFields, zap.Any(key, v))
		}
	}

	return &zapLogger{
		log: l.log.With(zapFields...),
	}
}
