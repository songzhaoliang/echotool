package binder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestXMLBodyBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`<user><id>1</id><name>peter</name></user>`))
	req.Header.Set("Content-Type", "application/xml")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	u := &User{}
	err := XMLBodyBinder.Bind(c, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "peter", u.Name)
}
