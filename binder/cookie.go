package binder

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

var CookieBinder = &cookieBinder{}

type cookieBinder struct {
}

var _ Binder = (*cookieBinder)(nil)

func (cookieBinder) Bind(c echo.Context, obj interface{}) error {
	return Bind(obj, parseCookie(c.Cookies()), TagCookie, false)
}

func parseCookie(cookies []*http.Cookie) (v url.Values) {
	v = make(url.Values)
	for _, cookie := range cookies {
		v[cookie.Name] = append(v[cookie.Name], cookie.Value)
	}
	return
}
