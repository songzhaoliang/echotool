package binder

import (
	"github.com/labstack/echo"
)

var FormPostBinder = &formPostBinder{}

type formPostBinder struct {
}

var _ Binder = (*formPostBinder)(nil)

func (formPostBinder) Bind(c echo.Context, obj interface{}) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	return bind(obj, c.Request().PostForm, TagForm, false)
}
