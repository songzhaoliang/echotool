package binder

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestFormMultipartBinder(t *testing.T) {
	body := &bytes.Buffer{}
	mv := multipart.NewWriter(body)
	mv.WriteField("id", "1")
	mv.WriteField("name", "peter")
	mv.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", mv.FormDataContentType())

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	u := &User{}
	err := FormMultipartBinder.Bind(c, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "peter", u.Name)
}
