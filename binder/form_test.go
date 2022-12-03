package binder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFormBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/?id=1", strings.NewReader("name=peter"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	u := &User{}
	err := FormBinder.Bind(c, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "peter", u.Name)
}
