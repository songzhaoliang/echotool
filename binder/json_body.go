package binder

import (
	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool/json"
)

var (
	EnableDecoderUseNumber             = false
	EnableDecoderDisallowUnknownFields = false
)

var JSONBodyBinder = &jsonBodyBinder{}

type jsonBodyBinder struct {
}

var _ Binder = (*jsonBodyBinder)(nil)

func (jsonBodyBinder) Bind(c echo.Context, obj interface{}) error {
	decoder := json.NewDecoder(c.Request().Body)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	return decoder.Decode(obj)
}
