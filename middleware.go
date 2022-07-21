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
			if handy.IsEmptyStr(id) {
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
