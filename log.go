package echotool

import (
	"io"
	"time"

	"github.com/popeyeio/handy"
	"github.com/songzhaoliang/echotool/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger = NewDefaultLogger()

func NewDefaultLogger() *zap.SugaredLogger {
	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
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
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	l, _ := cfg.Build()
	return l.Sugar()
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

func FasterJSONReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc
}
