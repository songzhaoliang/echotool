package echotool

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/songzhaoliang/echotool/metric"
	"github.com/songzhaoliang/echotool/pprof"
)

type HandlerFunc func(echo.Context, *Context)
type HandlerFuncsChain []HandlerFunc

type Engine struct {
	middlewares HandlerFuncsChain
	finisher    HandlerFunc
	aborter     HandlerFunc
	contextPool sync.Pool
}

type Option func(*Engine)

func WithFinisher(finisher HandlerFunc) Option {
	return func(e *Engine) {
		if finisher != nil {
			e.finisher = finisher
		}
	}
}

func WithAborter(aborter HandlerFunc) Option {
	return func(e *Engine) {
		if aborter != nil {
			e.aborter = aborter
		}
	}
}

func NewEngine(opts ...Option) *Engine {
	e := &Engine{
		finisher: GetCommonFinisher(),
		aborter:  GetCommonAborter(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func NewDefaultEcho() *echo.Echo {
	r := echo.New()
	r.Use(middleware.Recover())
	r.Use(SetRequestID(GetUUID))
	pprof.Register(r)
	metric.Register(r)
	return r
}

func NewDefaultEngine(opts ...Option) *Engine {
	e := NewEngine(opts...)
	e.Use(AddTraceID(GetRequestID))
	e.Use(PrintRequest())
	return e
}

func (e *Engine) Use(middlewares ...HandlerFunc) *Engine {
	e.middlewares = append(e.middlewares, middlewares...)
	return e
}

func (e *Engine) EchoHandler(handlers ...HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if len(handlers) == 0 {
			return nil
		}

		ec := e.acquireContext()
		defer e.releaseContext(ec)

		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(*EchotoolError); ok {
					ec.Abort(err.GetCode(), err.GetError())
					ReleaseEchotoolError(err)
					e.aborter(c, ec)
				} else {
					panic(r)
				}
			}
		}()

		ec.handlerName = GetHandlerName(handlers[0])
		ec.handlers = append(ec.handlers, e.middlewares...)
		ec.handlers = append(ec.handlers, handlers...)

		for _, handler := range ec.handlers {
			handler(c, ec)

			if !ec.IsOK() {
				break
			}
		}

		if ec.IsOK() {
			e.finisher(c, ec)
		} else {
			e.aborter(c, ec)
		}

		return nil
	}
}

func (e *Engine) acquireContext() (ec *Context) {
	if v := e.contextPool.Get(); v != nil {
		ec = v.(*Context)
	} else {
		ec = &Context{}
	}

	ec.engine = e
	ec.ok = true
	ec.customValues = make(map[string]string)
	ec.startTime = time.Now()
	return
}

func (e *Engine) releaseContext(ec *Context) {
	if ec != nil {
		ec.reset()
		e.contextPool.Put(ec)
	}
}
