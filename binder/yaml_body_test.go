package binder

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestYAMLBodyBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", encodeYAML(&User{1, "peter"}))
	req.Header.Set("Content-Type", "application/yaml")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	u := &User{}
	err := YAMLBodyBinder.Bind(c, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "peter", u.Name)
}

func encodeYAML(v interface{}) io.Reader {
	bs, _ := yaml.Marshal(v)
	return bytes.NewReader(bs)
}
