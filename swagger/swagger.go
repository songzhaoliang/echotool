package swagger

import (
	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool/util"
	es "github.com/swaggo/echo-swagger"
)

const (
	DefaultPrefix = "/swagger"
)

func Register(r *echo.Echo, prefixes ...string) {
	register(r.Group(util.GetPrefix(append(prefixes, DefaultPrefix)...)))
}

func RouterRegister(g *echo.Group, prefixes ...string) {
	register(g.Group(util.GetPrefix(append(prefixes, DefaultPrefix)...)))
}

func register(g *echo.Group) {
	g.GET("/*", es.WrapHandler)
}
