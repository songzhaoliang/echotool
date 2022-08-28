package binder

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int    `header:"X-Id" param:"id" form:"id" json:"id" xml:"id" msgpack:"id" yaml:"id" env:"ID"`
	Name string `header:"X-Name" param:"name" form:"name" json:"name" xml:"name"  msgpack:"name" yaml:"name" env:"NAME"`
}

func TestHeaderBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Id", "1")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	u := &User{}
	err := HeaderBinder.Bind(c, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
}
