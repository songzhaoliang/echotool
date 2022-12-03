package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func GetPrefix(prefixes ...string) (prefix string) {
	if len(prefixes) > 0 {
		prefix = prefixes[0]
	}
	return
}
