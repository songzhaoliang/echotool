package binder

import (
	"errors"
	"net/textproto"

	"github.com/labstack/echo/v4"
	"github.com/popeyeio/handy"
)

const (
	memoryMax = 1 << 25

	TagHeader   = "header"
	TagParam    = "param"
	TagForm     = "form"
	TagJSON     = "json"
	TagXML      = "xml"
	TagProtobuf = "protobuf"
	TagMsgpack  = "msgpack"
	TagYAML     = "yaml"
	TagEnv      = "env"
)

var (
	ErrInvalidType = errors.New("invalid type")
)

type Binder interface {
	Bind(echo.Context, interface{}) error
}

type BeforeBinder interface {
	BeforeBind(echo.Context) error
}

type AfterBinder interface {
	AfterBind(echo.Context) error
}

func canonicalKey(key string, canonical bool) string {
	if !canonical {
		return key
	}
	return textproto.CanonicalMIMEHeaderKey(key)
}

func convertValue(value string) string {
	if !handy.IsEmptyStr(value) {
		return value
	}
	return "0"
}
