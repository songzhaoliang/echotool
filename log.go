package echotool

import (
	"time"

	"github.com/popeyeio/handy"
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
			EncodeDuration: zapcore.StringDurationEncoder,
			// EncodeCaller:   zapcore.ShortCallerEncoder,
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

func buildLogger(ec *Context) (l *zap.SugaredLogger) {
	l = logger
	if v := ec.GetNamedValue(); !handy.IsEmptyStr(v) {
		l = l.Named(v)
	}

	for k, v := range ec.GetCustomValues() {
		l = l.With(zap.String(k, v))
	}
	return
}
