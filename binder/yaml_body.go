package binder

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v2"
)

var YAMLBodyBinder = &yamlBodyBinder{}

type yamlBodyBinder struct {
}

var _ Binder = (*yamlBodyBinder)(nil)

func (yamlBodyBinder) Bind(c echo.Context, obj interface{}) error {
	return yaml.NewDecoder(c.Request().Body).Decode(obj)
}
