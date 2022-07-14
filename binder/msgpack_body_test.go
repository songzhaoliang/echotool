package binder

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

func TestMsgpackBodyBinder(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", encodeMsgpack(&User{1, "peter"}))
	req.Header.Set("Content-Type", "application/msgpack")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	u := &User{}
	err := MsgpackBodyBinder.Bind(c, u)

	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "peter", u.Name)
}

func encodeMsgpack(v interface{}) *bytes.Buffer {
	buf := &bytes.Buffer{}
	codec.NewEncoder(buf, &codec.MsgpackHandle{}).Encode(v)
	return buf
}
