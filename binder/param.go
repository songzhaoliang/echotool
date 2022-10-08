package binder

import (
	"net/url"

	"github.com/labstack/echo"
)

var ParamBinder = &paramBinder{}

type paramBinder struct {
}

var _ Binder = (*paramBinder)(nil)

func (paramBinder) Bind(c echo.Context, obj interface{}) error {
	return Bind(obj, parseParam(c.ParamNames(), c.ParamValues()), TagParam, false)
}

func parseParam(names, values []string) (v url.Values) {
	v = make(url.Values)
	for i, name := range names {
		v[name] = append(v[name], values[i])
	}
	return
}
