package echotool

import (
	"fmt"

	"github.com/popeyeio/handy"
)

type Context struct {
	engine   *Engine
	handlers HandlerFuncsChain

	ok   bool
	code int
	data interface{}
	err  error
}

var _ fmt.Stringer = (*Context)(nil)

func (ec *Context) String() string {
	if ec.ok {
		return fmt.Sprintf("code:%d, data:%s", ec.code, handy.Stringify(ec.data))
	}
	return fmt.Sprintf("code:%d, error:%v", ec.code, ec.err)
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

func (ec *Context) reset() {
	ec.engine = nil
	ec.handlers = ec.handlers[:0]
	ec.ok = false
	ec.code = 0
	ec.data = nil
	ec.err = nil
}
