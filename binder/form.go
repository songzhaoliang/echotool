package binder

import (
	"github.com/labstack/echo"
)

var FormBinder = &formBinder{}

type formBinder struct {
}

var _ Binder = (*formBinder)(nil)

func (formBinder) Bind(c echo.Context, obj interface{}) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}

	c.Request().ParseMultipartForm(memoryMax)
	return bind(obj, c.Request().Form, TagForm, false)
}
