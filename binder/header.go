package binder

import (
	"github.com/labstack/echo"
)

var HeaderBinder = &headerBinder{}

type headerBinder struct {
}

var _ Binder = (*headerBinder)(nil)

func (headerBinder) Bind(c echo.Context, obj interface{}) error {
	return Bind(obj, c.Request().Header, TagHeader, true)
}
