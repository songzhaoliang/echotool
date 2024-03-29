package echotool

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/popeyeio/handy"
	"moul.io/http2curl"
)

const (
	KeyRequestID = "x-request-id"
)

func SetRequestID(f func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, _ := HeaderString(c, KeyRequestID)
			if handy.IsEmptyStr(id) && f != nil {
				id = f(c)
			}

			c.Set(KeyRequestID, id)

			return next(c)
		}
	}
}

func GetRequestID(c echo.Context) string {
	val := c.Get(KeyRequestID)
	if val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return handy.StrEmpty
}

// AddTraceID adds trace id to log for troubleshooting problems and so on.
func AddTraceID(f func(c echo.Context) string) HandlerFunc {
	return func(c echo.Context, ec *Context) {
		if f != nil {
			ec.SetNamedValue(f(c))
		}
	}
}

// AddNotice adds kv pair to log for troubleshooting problems and so on.
func AddNotice(key, value string) HandlerFunc {
	return func(c echo.Context, ec *Context) {
		if !handy.IsEmptyStr(key) || !handy.IsEmptyStr(value) {
			ec.SetCustomValue(key, value)
		}
	}
}

func PrintRequest() HandlerFunc {
	return func(c echo.Context, ec *Context) {
		if cmd, err := http2curl.GetCurlCommand(c.Request()); err == nil {
			CtxInfoKV(ec, cmd.String())
		}
	}
}

func DefaultCORS() echo.MiddlewareFunc {
	return middleware.CORS()
}

func UnsafeCORS() echo.MiddlewareFunc {
	middleware.DefaultCORSConfig.AllowCredentials = true
	middleware.DefaultCORSConfig.UnsafeWildcardOriginWithAllowCredentials = true // send origin back
	middleware.DefaultCORSConfig.MaxAge = 2592000                                // the result of preflight can be cached one month
	middleware.DefaultCORSConfig.ExposeHeaders = []string{echo.HeaderVary}       // headers can be accessed by clients
	return DefaultCORS()
}
