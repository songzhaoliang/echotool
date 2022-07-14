package binder

import (
	"github.com/labstack/echo"
)

var FormMultipartBinder = &formMultipartBinder{}

type formMultipartBinder struct {
}

var _ Binder = (*formMultipartBinder)(nil)

func (formMultipartBinder) Bind(c echo.Context, obj interface{}) error {
	if err := c.Request().ParseMultipartForm(memoryMax); err != nil {
		return err
	}
	return bind(obj, c.Request().MultipartForm.Value, TagForm, false)
}
