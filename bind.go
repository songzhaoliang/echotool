package echotool

import (
	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool/binder"
	"github.com/songzhaoliang/echotool/validator"
)

// BindHeader needs tag "header" in fields of v.
// The value of tag "header" is automatically converted to the canonical format.
func BindHeader(c echo.Context, v interface{}) error {
	return binder.HeaderBinder.Bind(c, v)
}

func MustBindHeader(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, BindHeader(c, v)
	}, CodeBindErr, cbs...)
}

// BindParam needs tag "param" in fields of v.
func BindParam(c echo.Context, v interface{}) error {
	return binder.ParamBinder.Bind(c, v)
}

func MustBindParam(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, BindParam(c, v)
	}, CodeBindErr, cbs...)
}

// FormBindQuery needs tag "form" in fields of v.
func FormBindQuery(c echo.Context, v interface{}) error {
	return binder.QueryBinder.Bind(c, v)
}

func MustFormBindQuery(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, FormBindQuery(c, v)
	}, CodeBindErr, cbs...)
}

// FormBindBody needs tag "form" in fields of v.
func FormBindBody(c echo.Context, v interface{}) error {
	return binder.FormPostBinder.Bind(c, v)
}

func MustFormBindBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, FormBindBody(c, v)
	}, CodeBindErr, cbs...)
}

// FormBindQueryBody needs tag "form" in fields of v.
func FormBindQueryBody(c echo.Context, v interface{}) error {
	return binder.FormBinder.Bind(c, v)
}

func MustFormBindQueryBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, FormBindQueryBody(c, v)
	}, CodeBindErr, cbs...)
}

// FormBindMultipart needs tag "form" in fields of v.
func FormBindMultipart(c echo.Context, v interface{}) error {
	return binder.FormMultipartBinder.Bind(c, v)
}

func MustFormBindMultipart(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, FormBindMultipart(c, v)
	}, CodeBindErr, cbs...)
}

// JSONBindBody needs tag "json" in fields of v.
func JSONBindBody(c echo.Context, v interface{}) error {
	return binder.JSONBodyBinder.Bind(c, v)
}

func MustJSONBindBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, JSONBindBody(c, v)
	}, CodeBindErr, cbs...)
}

// XMLBindBody needs tag "xml" in fields of v.
func XMLBindBody(c echo.Context, v interface{}) error {
	return binder.XMLBodyBinder.Bind(c, v)
}

func MustXMLBindBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, XMLBindBody(c, v)
	}, CodeBindErr, cbs...)
}

// ProtobufBindBody needs tag "protobuf" in fields of v.
func ProtobufBindBody(c echo.Context, v interface{}) error {
	return binder.ProtobufBodyBinder.Bind(c, v)
}

func MustProtobufBindBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, ProtobufBindBody(c, v)
	}, CodeBindErr, cbs...)
}

// MsgpackBindBody needs tag "msgpack" in fields of v.
func MsgpackBindBody(c echo.Context, v interface{}) error {
	return binder.MsgpackBodyBinder.Bind(c, v)
}

func MustMsgpackBindBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, MsgpackBindBody(c, v)
	}, CodeBindErr, cbs...)
}

// YAMLBindBody needs tag "yaml" in fields of v.
func YAMLBindBody(c echo.Context, v interface{}) error {
	return binder.YAMLBodyBinder.Bind(c, v)
}

func MustYAMLBindBody(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, YAMLBindBody(c, v)
	}, CodeBindErr, cbs...)
}

// BindEnv needs tag "env" in fields of v.
func BindEnv(c echo.Context, v interface{}) error {
	return binder.EnvBinder.Bind(c, v)
}

func MustBindEnv(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, BindEnv(c, v)
	}, CodeBindErr, cbs...)
}

// BindCookie needs tag "cookie" in fields of v.
func BindCookie(c echo.Context, v interface{}) error {
	return binder.CookieBinder.Bind(c, v)
}

func MustBindCookie(c echo.Context, v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, BindCookie(c, v)
	}, CodeBindErr, cbs...)
}

// Validate needs tag "valid" in fields of v.
func Validate(v interface{}) error {
	return validator.EchotoolValidator.ValidateStruct(v)
}

func MustValidate(v interface{}, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, Validate(v)
	}, CodeValidateErr, cbs...)
}

const (
	BValidator = 1 << iota
	BHeader
	BParam
	BFormQuery
	BFormBody
	BFormQueryBody
	BFormMultipart
	BJSONBody
	BXMLBody
	BProtobufBody
	BMsgpackBody
	BYAMLBody
	BEnv
	BCookie
)

var funcs = map[int]func(echo.Context, interface{}) error{
	BHeader:        BindHeader,
	BParam:         BindParam,
	BFormQuery:     FormBindQuery,
	BFormBody:      FormBindBody,
	BFormQueryBody: FormBindQueryBody,
	BFormMultipart: FormBindMultipart,
	BJSONBody:      JSONBindBody,
	BXMLBody:       XMLBindBody,
	BProtobufBody:  ProtobufBindBody,
	BMsgpackBody:   MsgpackBindBody,
	BYAMLBody:      YAMLBindBody,
	BEnv:           BindEnv,
	BCookie:        BindCookie,
}

func RegisterBinder(flag int, fn func(echo.Context, interface{}) error) bool {
	if _, exists := funcs[flag]; exists {
		return false
	}

	funcs[flag] = fn
	return true
}

func ForceRegisterBinder(flag int, fn func(echo.Context, interface{}) error) {
	funcs[flag] = fn
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
			return AcquireEchotoolError(CodeValidateErr, err)
		}
	}

	if obj, ok := v.(binder.AfterBinder); ok {
		if err = obj.AfterBind(c); err != nil {
			return
		}
	}
	return
}

func MustBind(c echo.Context, v interface{}, flag int, cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, Bind(c, v, flag)
	}, CodeBindErr, cbs...)
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

func (p *proxy) FormBindMultipart() *proxy {
	p.flag |= BFormMultipart
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

func (p *proxy) BindEnv() *proxy {
	p.flag |= BEnv
	return p
}

func (p *proxy) BindCookie() *proxy {
	p.flag |= BCookie
	return p
}

func (p *proxy) Validate() *proxy {
	p.flag |= BValidator
	return p
}

func (p *proxy) End() error {
	return Bind(p.c, p.v, p.flag)
}

func (p *proxy) MustEnd(cbs ...CallbackFunc) {
	MustDoCallback(func() (interface{}, error) {
		return nil, p.End()
	}, CodeBindErr, cbs...)
}
