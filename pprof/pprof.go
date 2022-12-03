package pprof

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool/util"
)

const (
	DefaultPrefix = "/debug/pprof"
)

func Register(r *echo.Echo, prefixes ...string) {
	register(r.Group(util.GetPrefix(append(prefixes, DefaultPrefix)...)))
}

func RouterRegister(g *echo.Group, prefixes ...string) {
	register(g.Group(util.GetPrefix(append(prefixes, DefaultPrefix)...)))
}

func register(g *echo.Group) {
	g.GET("/", util.WrapF(pprof.Index))
	g.GET("/cmdline", util.WrapF(pprof.Cmdline))
	g.GET("/profile", util.WrapF(pprof.Profile))
	g.POST("/symbol", util.WrapF(pprof.Symbol))
	g.GET("/symbol", util.WrapF(pprof.Symbol))
	g.GET("/trace", util.WrapF(pprof.Trace))
	g.GET("/allocs", util.WrapH(pprof.Handler("allocs")))
	g.GET("/block", util.WrapH(pprof.Handler("block")))
	g.GET("/goroutine", util.WrapH(pprof.Handler("goroutine")))
	g.GET("/heap", util.WrapH(pprof.Handler("heap")))
	g.GET("/mutex", util.WrapH(pprof.Handler("mutex")))
	g.GET("/threadcreate", util.WrapH(pprof.Handler("threadcreate")))
}
