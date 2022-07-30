package echotool

import (
	"github.com/labstack/echo"
	"github.com/popeyeio/handy"
)

const (
	KeyRequestID = "x_request_id"
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
		if key != handy.StrEmpty || value != handy.StrEmpty {
			ec.SetCustomValue(key, value)
		}
	}
}
