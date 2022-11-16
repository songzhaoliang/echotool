package echotool

import (
	"context"
	"fmt"
	"time"

	"github.com/popeyeio/handy"
)

type Context struct {
	engine      *Engine
	handlers    HandlerFuncsChain
	handlerName string

	ok   bool
	code int
	data interface{}
	err  error

	namedValue   string
	customValues map[string]string
}

var _ context.Context = (*Context)(nil)
var _ fmt.Stringer = (*Context)(nil)

func (ec *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ec *Context) Done() <-chan struct{} {
	return nil
}

func (ec *Context) Err() error {
	return nil
}

func (ec *Context) Value(key interface{}) interface{} {
	return nil
}

func (ec *Context) String() string {
	if ec.ok {
		return fmt.Sprintf("code:%d, data:%s", ec.code, handy.Stringify(ec.data))
	}
	return fmt.Sprintf("code:%d, error:%v", ec.code, ec.err)
}

func (ec *Context) GetHandlerName() string {
	return ec.handlerName
}

func (ec *Context) SetHandlerName(name string) {
	ec.handlerName = name
}

func (ec *Context) Finish(code int, data interface{}) {
	ec.ok = true
	ec.code = code
	ec.data = data
}

func (ec *Context) Abort(code int, err error) {
	ec.ok = false
	ec.code = code
	ec.err = err
}

func (ec *Context) IsOK() bool {
	return ec.ok
}

func (ec *Context) GetCode() int {
	return ec.code
}

func (ec *Context) GetData() interface{} {
	return ec.data
}

func (ec *Context) GetError() error {
	return ec.err
}

func (ec *Context) GetNamedValue() string {
	return ec.namedValue
}

func (ec *Context) SetNamedValue(value string) {
	ec.namedValue = value
}

func (ec *Context) GetCustomValue(key string) (value string, exists bool) {
	value, exists = ec.customValues[key]
	return
}

func (ec *Context) SetCustomValue(key, value string) {
	ec.customValues[key] = value
}

func (ec *Context) GetCustomValues() map[string]string {
	return ec.customValues
}

func (ec *Context) reset() {
	ec.engine = nil
	ec.handlers = ec.handlers[:0]
	ec.handlerName = handy.StrEmpty
	ec.ok = false
	ec.code = 0
	ec.data = nil
	ec.err = nil
	ec.namedValue = handy.StrEmpty
	ec.customValues = nil
}
