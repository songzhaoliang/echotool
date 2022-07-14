package binder

import (
	"encoding/xml"

	"github.com/labstack/echo"
)

var XMLBodyBinder = &xmlBodyBinder{}

type xmlBodyBinder struct {
}

var _ Binder = (*xmlBodyBinder)(nil)

func (xmlBodyBinder) Bind(c echo.Context, obj interface{}) error {
	return xml.NewDecoder(c.Request().Body).Decode(obj)
}
