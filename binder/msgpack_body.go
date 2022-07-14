package binder

import (
	"github.com/labstack/echo"
	"github.com/ugorji/go/codec"
)

var MsgpackBodyBinder = &msgpackBodyBinder{}

type msgpackBodyBinder struct {
}

var _ Binder = (*msgpackBodyBinder)(nil)

func (msgpackBodyBinder) Bind(c echo.Context, obj interface{}) error {
	cdc := &codec.MsgpackHandle{}
	return codec.NewDecoder(c.Request().Body, cdc).Decode(&obj)
}
