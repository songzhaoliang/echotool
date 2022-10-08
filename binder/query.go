package binder

import (
	"github.com/labstack/echo"
)

var QueryBinder = &queryBinder{}

type queryBinder struct {
}

var _ Binder = (*queryBinder)(nil)

func (queryBinder) Bind(c echo.Context, obj interface{}) error {
	return Bind(obj, c.Request().URL.Query(), TagForm, false)
}
