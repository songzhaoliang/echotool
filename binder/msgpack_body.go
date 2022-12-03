package binder

import (
	"github.com/labstack/echo/v4"
	"github.com/ugorji/go/codec"
)

var MsgpackBodyBinder = &msgpackBodyBinder{}

type msgpackBodyBinder struct {
}

var _ Binder = (*msgpackBodyBinder)(nil)

func (msgpackBodyBinder) Bind(c echo.Context, obj interface{}) error {
	return codec.NewDecoder(c.Request().Body, &codec.MsgpackHandle{}).Decode(&obj)
}
