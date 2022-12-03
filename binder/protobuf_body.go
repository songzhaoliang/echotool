package binder

import (
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo/v4"
)

var ProtobufBodyBinder = &protobufBodyBinder{}

type protobufBodyBinder struct {
}

var _ Binder = (*protobufBodyBinder)(nil)

func (protobufBodyBinder) Bind(c echo.Context, obj interface{}) error {
	bs, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	return proto.Unmarshal(bs, obj.(proto.Message))
}
