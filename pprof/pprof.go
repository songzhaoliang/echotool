package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/labstack/echo"
)

const (
	DefaultPrefix = "/debug/pprof"
)

func Register(r *echo.Echo, prefixes ...string) {
	register(r.Group(getPrefix(prefixes...)))
}

func RouterRegister(g *echo.Group, prefixes ...string) {
	register(g.Group(getPrefix(prefixes...)))
}

func register(g *echo.Group) {
	g.GET("/", WrapF(pprof.Index))
	g.GET("/cmdline", WrapF(pprof.Cmdline))
	g.GET("/profile", WrapF(pprof.Profile))
	g.POST("/symbol", WrapF(pprof.Symbol))
	g.GET("/symbol", WrapF(pprof.Symbol))
	g.GET("/trace", WrapF(pprof.Trace))
	g.GET("/allocs", WrapH(pprof.Handler("allocs")))
	g.GET("/block", WrapH(pprof.Handler("block")))
	g.GET("/goroutine", WrapH(pprof.Handler("goroutine")))
	g.GET("/heap", WrapH(pprof.Handler("heap")))
	g.GET("/mutex", WrapH(pprof.Handler("mutex")))
	g.GET("/threadcreate", WrapH(pprof.Handler("threadcreate")))
}

func getPrefix(prefixes ...string) string {
	prefix := DefaultPrefix
	if len(prefixes) > 0 {
		prefix = prefixes[0]
	}
	return prefix
}

func WrapF(f http.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		f(c.Response().Writer, c.Request())
		return nil
	}
}

func WrapH(h http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
