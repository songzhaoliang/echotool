package echotool

import (
	"context"
	"net"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap/zapcore"
)

type RedisLogger struct {
	Level zapcore.Level
}

var _ redis.Hook = (*RedisLogger)(nil)

func NewRedisLogger(level string) *RedisLogger {
	lv, err := zapcore.ParseLevel(level)
	if err != nil {
		lv = zapcore.DebugLevel
	}

	return &RedisLogger{
		Level: lv,
	}
}

func (l *RedisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		ec, _ := ctx.(*Context)
		CtxPrint(ec, l.Level, "[redis] dial network %s, addr %s", network, addr)

		return next(ctx, network, addr)
	}
}

func (l *RedisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		ec, _ := ctx.(*Context)
		CtxPrint(ec, l.Level, "[redis] %s", cmd.String())

		return next(ctx, cmd)
	}
}

func (l *RedisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		ec, _ := ctx.(*Context)
		for _, cmd := range cmds {
			CtxPrint(ec, l.Level, "[redis] %s", cmd.String())
		}

		return next(ctx, cmds)
	}
}
