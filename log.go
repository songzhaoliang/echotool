package echotool

import (
	"context"
	"io"
	"os"
	"time"

	rl "github.com/lestrrat-go/file-rotatelogs"
	"github.com/popeyeio/handy"
	"github.com/songzhaoliang/echotool/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

var logger = NewDefaultLogger()

func NewDefaultLogger() *zap.SugaredLogger {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "console",
		EncoderConfig:    NewDefaultEncodeConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	l, _ := cfg.Build()
	return l.Sugar()
}

// RotateConfig is the config for rotating logs.
// The format of suffix refers to https://github.com/lestrrat-go/strftime.
type RotateConfig struct {
	EncoderConfig zapcore.EncoderConfig
	Paths         []string
	Suffix        string
	Level         zapcore.Level
	RotateTime    time.Duration
	TTL           time.Duration
}

type RotateConfigOption func(*RotateConfig)

func WithEncoderConfig(cfg zapcore.EncoderConfig) RotateConfigOption {
	return func(c *RotateConfig) {
		c.EncoderConfig = cfg
	}
}

func WithPaths(paths []string) RotateConfigOption {
	return func(c *RotateConfig) {
		if len(paths) > 0 {
			c.Paths = paths
		}
	}
}

func WithSuffix(suffix string) RotateConfigOption {
	return func(c *RotateConfig) {
		if !handy.IsEmptyStr(suffix) {
			c.Suffix = suffix
		}
	}
}

func WithLevel(level zapcore.Level) RotateConfigOption {
	return func(c *RotateConfig) {
		c.Level = level
	}
}

func WithRotateTime(t time.Duration) RotateConfigOption {
	return func(c *RotateConfig) {
		if t > 0 {
			c.RotateTime = t
		}
	}
}

func WithTTL(ttl time.Duration) RotateConfigOption {
	return func(c *RotateConfig) {
		if ttl > 0 {
			c.TTL = ttl
		}
	}
}

func NewRotateLogger(opts ...RotateConfigOption) (*zap.SugaredLogger, error) {
	cfg := &RotateConfig{
		EncoderConfig: NewDefaultEncodeConfig(),
		Paths:         []string{"stdout"},
		Suffix:        "%Y%m%d_%H",
		Level:         zapcore.DebugLevel,
		RotateTime:    time.Hour,
		TTL:           time.Hour * 24 * 7,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	level := zap.NewAtomicLevelAt(cfg.Level)
	enc := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	var cores []zapcore.Core
	for _, path := range cfg.Paths {
		switch path {
		case "stdout":
			cores = append(cores, zapcore.NewCore(enc, os.Stdout, level))
		case "stderr":
			cores = append(cores, zapcore.NewCore(enc, os.Stderr, level))
		default:
			w, err := rl.New(path+"."+cfg.Suffix,
				rl.WithLinkName(path),
				rl.WithRotationTime(cfg.RotateTime),
				rl.WithMaxAge(cfg.TTL))
			if err != nil {
				return nil, err
			}

			cores = append(cores, zapcore.NewCore(enc, zapcore.AddSync(w), level))
		}
	}

	return zap.New(zapcore.NewTee(cores...)).Sugar(), nil
}

func NewDefaultEncodeConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:     "T",
		LevelKey:    "L",
		NameKey:     "N",
		CallerKey:   "C",
		MessageKey:  "M",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration:      zapcore.StringDurationEncoder,
		NewReflectedEncoder: FasterJSONReflectedEncoder,
	}
}

func FasterJSONReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc
}

func SetLogger(l *zap.SugaredLogger) {
	if l != nil {
		logger = l
	}
}

// FlushLog flushes any buffered log entries.
func FlushLog() error {
	return logger.Sync()
}

func CtxDebug(ec *Context, format string, args ...interface{}) {
	buildLogger(ec).Debugf(format, args...)
}

func CtxInfo(ec *Context, format string, args ...interface{}) {
	buildLogger(ec).Infof(format, args...)
}

func CtxWarn(ec *Context, format string, args ...interface{}) {
	buildLogger(ec).Warnf(format, args...)
}

func CtxError(ec *Context, format string, args ...interface{}) {
	buildLogger(ec).Errorf(format, args...)
}

// CtxFatal logs a message, then calls os.Exit.
func CtxFatal(ec *Context, format string, args ...interface{}) {
	buildLogger(ec).Fatalf(format, args...)
}

// CtxPanic logs a message, then panics.
func CtxPanic(ec *Context, format string, args ...interface{}) {
	buildLogger(ec).Panicf(format, args...)
}

func CtxDebugKV(ec *Context, msg string, fields ...zap.Field) {
	buildLogger(ec).Desugar().Debug(msg, fields...)
}

func CtxInfoKV(ec *Context, msg string, fields ...zap.Field) {
	buildLogger(ec).Desugar().Info(msg, fields...)
}

func CtxWarnKV(ec *Context, msg string, fields ...zap.Field) {
	buildLogger(ec).Desugar().Warn(msg, fields...)
}

func CtxErrorKV(ec *Context, msg string, fields ...zap.Field) {
	buildLogger(ec).Desugar().Error(msg, fields...)
}

// CtxFatalKV logs a message, then calls os.Exit.
func CtxFatalKV(ec *Context, msg string, fields ...zap.Field) {
	buildLogger(ec).Desugar().Fatal(msg, fields...)
}

// CtxPanicKV logs a message, then panics.
func CtxPanicKV(ec *Context, msg string, fields ...zap.Field) {
	buildLogger(ec).Desugar().Panic(msg, fields...)
}

func buildLogger(ec *Context) (l *zap.SugaredLogger) {
	l = logger
	if ec == nil {
		return
	}

	if v := ec.GetNamedValue(); !handy.IsEmptyStr(v) {
		l = l.Named(v)
	}

	for k, v := range ec.GetCustomValues() {
		l = l.With(zap.String(k, v))
	}
	return
}

func Debug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal logs a message, then calls os.Exit.
func Fatal(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panic logs a message, then panics.
func Panic(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func DebugKV(msg string, fields ...zap.Field) {
	logger.Desugar().Debug(msg, fields...)
}

func InfoKV(msg string, fields ...zap.Field) {
	logger.Desugar().Info(msg, fields...)
}

func WarnKV(msg string, fields ...zap.Field) {
	logger.Desugar().Warn(msg, fields...)
}

func ErrorKV(msg string, fields ...zap.Field) {
	logger.Desugar().Error(msg, fields...)
}

// FatalKV logs a message, then calls os.Exit.
func FatalKV(msg string, fields ...zap.Field) {
	logger.Desugar().Fatal(msg, fields...)
}

// PanicKV logs a message, then panics.
func PanicKV(msg string, fields ...zap.Field) {
	logger.Desugar().Panic(msg, fields...)
}

type GORMLogger struct {
	SlowTime       time.Duration
	IgnoreNoRecord bool
}

var _ gl.Interface = (*GORMLogger)(nil)

func NewDefaultGORMLogger() *GORMLogger {
	return &GORMLogger{
		SlowTime:       time.Millisecond * 200,
		IgnoreNoRecord: true,
	}
}

func (l *GORMLogger) LogMode(gl.LogLevel) gl.Interface {
	return l
}

func (l *GORMLogger) Info(ctx context.Context, format string, args ...interface{}) {
	ec, _ := ctx.(*Context)
	CtxInfo(ec, format, args...)
}

func (l *GORMLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	ec, _ := ctx.(*Context)
	CtxWarn(ec, format, args...)
}

func (l *GORMLogger) Error(ctx context.Context, format string, args ...interface{}) {
	ec, _ := ctx.(*Context)
	CtxError(ec, format, args...)
}

func (l *GORMLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	ec, _ := ctx.(*Context)
	tns := time.Since(begin)
	tms := float64(tns.Nanoseconds()) / 1e6
	sql, rows := f()
	switch {
	case err != nil && (err != gorm.ErrRecordNotFound || !l.IgnoreNoRecord):
		CtxError(ec, "[%.3fms] [rows:%d] %s - %v", tms, rows, sql, err)
	case l.SlowTime > 0 && tns > l.SlowTime:
		CtxWarn(ec, "[%.3fms] [rows:%d] %s - slow sql > %v", tms, rows, sql, l.SlowTime)
	default:
		CtxInfo(ec, "[%.3fms] [rows:%d] %s", tms, rows, sql)
	}
}
