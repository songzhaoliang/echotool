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
	return bind(obj, parse(c.ParamNames(), c.ParamValues()), TagParam, false)
}

func parse(names, values []string) (v url.Values) {
	v = make(url.Values)
	for i, name := range names {
		if _, exists := v[name]; exists {
			v[name] = append(v[name], values[i])
		} else {
			v[name] = []string{values[i]}
		}
	}
	return
}
