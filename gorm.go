package echotool

import (
	"context"
	"time"

	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

type GORMLogger struct {
	Prefix         string
	SlowTime       time.Duration
	IgnoreNoRecord bool
}

var _ gl.Interface = (*GORMLogger)(nil)

func NewDefaultGORMLogger() *GORMLogger {
	return &GORMLogger{
		Prefix:         "[gorm]",
		SlowTime:       time.Millisecond * 200,
		IgnoreNoRecord: true,
	}
}

func (l *GORMLogger) LogMode(level gl.LogLevel) gl.Interface {
	return l
}

func (l *GORMLogger) Info(ctx context.Context, format string, args ...interface{}) {
	ec, _ := ctx.(*Context)
	CtxInfo(ec, l.addPrefix(format), args...)
}

func (l *GORMLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	ec, _ := ctx.(*Context)
	CtxWarn(ec, l.addPrefix(format), args...)
}

func (l *GORMLogger) Error(ctx context.Context, format string, args ...interface{}) {
	ec, _ := ctx.(*Context)
	CtxError(ec, l.addPrefix(format), args...)
}

func (l *GORMLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	tns := time.Since(begin)
	tms := float64(tns.Nanoseconds()) / 1e6
	sql, rows := f()
	switch {
	case err != nil && (err != gorm.ErrRecordNotFound || !l.IgnoreNoRecord):
		l.Error(ctx, "[%.3fms] [rows:%d] %s - %v", tms, rows, sql, err)
	case l.SlowTime > 0 && tns > l.SlowTime:
		l.Warn(ctx, "[%.3fms] [rows:%d] %s - slow sql > %v", tms, rows, sql, l.SlowTime)
	default:
		l.Info(ctx, "[%.3fms] [rows:%d] %s", tms, rows, sql)
	}
}

func (l *GORMLogger) addPrefix(msg string) string {
	return l.Prefix + " " + msg
}

func NewDefaultGORMConfig() *gorm.Config {
	return &gorm.Config{
		PrepareStmt: true,
		Logger:      NewDefaultGORMLogger(),
	}
}

func AutoIncrID(id *int64) {
	*id = 0
}
