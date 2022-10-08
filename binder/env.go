package binder

import (
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo"
)

var EnvBinder = &envBinder{}

type envBinder struct {
}

var _ Binder = (*envBinder)(nil)

func (envBinder) Bind(c echo.Context, obj interface{}) error {
	return Bind(obj, parseEnv(os.Environ()), TagEnv, false)
}

func parseEnv(envs []string) (v url.Values) {
	v = make(url.Values)
	for _, e := range envs {
		tokens := strings.SplitN(e, "=", 2)
		v[tokens[0]] = append(v[tokens[0]], tokens[1])
	}
	return
}
