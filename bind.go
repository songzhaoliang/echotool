package echotool

import (
	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool/binder"
	"github.com/songzhaoliang/echotool/validator"
)

// BindHeader needs tag "header" in fields of v.
// The value of tag "header" is automatically converted to the canonical format.
func BindHeader(c echo.Context, v interface{}) error {
	return binder.HeaderBinder.Bind(c, v)
}

// BindParam needs tag "param" in fields of v.
func BindParam(c echo.Context, v interface{}) error {
	return binder.ParamBinder.Bind(c, v)
}

// FormBindQuery needs tag "form" in fields of v.
func FormBindQuery(c echo.Context, v interface{}) error {
	return binder.QueryBinder.Bind(c, v)
}

// FormBindBody needs tag "form" in fields of v.
func FormBindBody(c echo.Context, v interface{}) error {
	return binder.FormPostBinder.Bind(c, v)
}

// FormBindQueryBody needs tag "form" in fields of v.
func FormBindQueryBody(c echo.Context, v interface{}) error {
	return binder.FormBinder.Bind(c, v)
}

// JSONBindBody needs tag "json" in fields of v.
func JSONBindBody(c echo.Context, v interface{}) error {
	return binder.JSONBodyBinder.Bind(c, v)
}

// XMLBindBody needs tag "xml" in fields of v.
func XMLBindBody(c echo.Context, v interface{}) error {
	return binder.XMLBodyBinder.Bind(c, v)
}

// ProtobufBindBody needs tag "protobuf" in fields of v.
func ProtobufBindBody(c echo.Context, v interface{}) error {
	return binder.ProtobufBodyBinder.Bind(c, v)
}

// MsgpackBindBody needs tag "msgpack" in fields of v.
func MsgpackBindBody(c echo.Context, v interface{}) error {
	return binder.MsgpackBodyBinder.Bind(c, v)
}

// YAMLBindBody needs tag "yaml" in fields of v.
func YAMLBindBody(c echo.Context, v interface{}) error {
	return binder.YAMLBodyBinder.Bind(c, v)
}

// Validate needs tag "valid" in fields of v.
func Validate(v interface{}) error {
	return validator.EchotoolValidator.ValidateStruct(v)
}

const (
	BValidator = 1 << iota
	BHeader
	BParam
	BFormQuery
	BFormBody
	BFormQueryBody
	BJSONBody
	BXMLBody
	BProtobufBody
	BMsgpackBody
	BYAMLBody
)

var funcs = map[int]func(echo.Context, interface{}) error{
	BHeader:        BindHeader,
	BParam:         BindParam,
	BFormQuery:     FormBindQuery,
	BFormBody:      FormBindBody,
	BFormQueryBody: FormBindQueryBody,
	BJSONBody:      JSONBindBody,
	BXMLBody:       XMLBindBody,
	BProtobufBody:  ProtobufBindBody,
	BMsgpackBody:   MsgpackBindBody,
	BYAMLBody:      YAMLBindBody,
}

func Bind(c echo.Context, v interface{}, flag int) (err error) {
	if obj, ok := v.(binder.BeforeBinder); ok {
		if err = obj.BeforeBind(c); err != nil {
			return
		}
	}

	for k, f := range funcs {
		if flag&k != 0 {
			if err = f(c, v); err != nil {
				return
			}
		}
	}

	// this must be the last one.
	if flag&BValidator != 0 {
		if err = Validate(v); err != nil {
			return
		}
	}

	if obj, ok := v.(binder.AfterBinder); ok {
		if err = obj.AfterBind(c); err != nil {
			return
		}
	}
	return
}

type proxy struct {
	c    echo.Context
	v    interface{}
	flag int
}

func New(c echo.Context, v interface{}) *proxy {
	return &proxy{
		c: c,
		v: v,
	}
}

func (p *proxy) BindHeader() *proxy {
	p.flag |= BHeader
	return p
}

func (p *proxy) BindParam() *proxy {
	p.flag |= BParam
	return p
}

func (p *proxy) FormBindQuery() *proxy {
	p.flag |= BFormQuery
	return p
}

func (p *proxy) FormBindBody() *proxy {
	p.flag |= BFormBody
	return p
}

func (p *proxy) FormBindQueryBody() *proxy {
	p.flag |= BFormQueryBody
	return p
}

func (p *proxy) JSONBindBody() *proxy {
	p.flag |= BJSONBody
	return p
}

func (p *proxy) XMLBindBody() *proxy {
	p.flag |= BXMLBody
	return p
}

func (p *proxy) ProtobufBindBody() *proxy {
	p.flag |= BProtobufBody
	return p
}

func (p *proxy) MsgpackBindBody() *proxy {
	p.flag |= BMsgpackBody
	return p
}

func (p *proxy) YAMLBindBody() *proxy {
	p.flag |= BYAMLBody
	return p
}

func (p *proxy) Validate() *proxy {
	p.flag |= BValidator
	return p
}

func (p *proxy) End() error {
	return Bind(p.c, p.v, p.flag)
}
