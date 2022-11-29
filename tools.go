package echotool

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/labstack/echo"
	"github.com/popeyeio/handy"
	"github.com/songzhaoliang/echotool/json"
)

type RunFunc func() (interface{}, error)
type CallbackFunc func(error)

func HeaderBool(c echo.Context, key string) (bool, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return false, NewHeaderEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func MustHeaderBool(c echo.Context, key string, cbs ...CallbackFunc) bool {
	result := MustDoCallback(func() (interface{}, error) {
		return HeaderBool(c, key)
	}, CodeBadRequest, cbs...)
	return result.(bool)
}

func HeaderInt64(c echo.Context, key string) (int64, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return 0, NewHeaderEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func MustHeaderInt64(c echo.Context, key string, cbs ...CallbackFunc) int64 {
	result := MustDoCallback(func() (interface{}, error) {
		return HeaderInt64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(int64)
}

func HeaderUint64(c echo.Context, key string) (uint64, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return 0, NewHeaderEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func MustHeaderUInt64(c echo.Context, key string, cbs ...CallbackFunc) uint64 {
	result := MustDoCallback(func() (interface{}, error) {
		return HeaderUint64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(uint64)
}

func HeaderString(c echo.Context, key string) (string, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewHeaderEmptyError(key)
	} else {
		return v, nil
	}
}

func MustHeaderString(c echo.Context, key string, cbs ...CallbackFunc) string {
	result := MustDoCallback(func() (interface{}, error) {
		return HeaderString(c, key)
	}, CodeBadRequest, cbs...)
	return result.(string)
}

func ParamBool(c echo.Context, key string) (bool, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return false, NewParamEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func MustParamBool(c echo.Context, key string, cbs ...CallbackFunc) bool {
	result := MustDoCallback(func() (interface{}, error) {
		return ParamBool(c, key)
	}, CodeBadRequest, cbs...)
	return result.(bool)
}

func ParamInt64(c echo.Context, key string) (int64, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return 0, NewParamEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func MustParamInt64(c echo.Context, key string, cbs ...CallbackFunc) int64 {
	result := MustDoCallback(func() (interface{}, error) {
		return ParamInt64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(int64)
}

func ParamUint64(c echo.Context, key string) (uint64, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return 0, NewParamEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func MustParamUInt64(c echo.Context, key string, cbs ...CallbackFunc) uint64 {
	result := MustDoCallback(func() (interface{}, error) {
		return ParamUint64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(uint64)
}

func ParamString(c echo.Context, key string) (string, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewParamEmptyError(key)
	} else {
		return v, nil
	}
}

func MustParamString(c echo.Context, key string, cbs ...CallbackFunc) string {
	result := MustDoCallback(func() (interface{}, error) {
		return ParamString(c, key)
	}, CodeBadRequest, cbs...)
	return result.(string)
}

func QueryBool(c echo.Context, key string) (bool, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return false, NewQueryEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func MustQueryBool(c echo.Context, key string, cbs ...CallbackFunc) bool {
	result := MustDoCallback(func() (interface{}, error) {
		return QueryBool(c, key)
	}, CodeBadRequest, cbs...)
	return result.(bool)
}

func QueryInt64(c echo.Context, key string) (int64, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return 0, NewQueryEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func MustQueryInt64(c echo.Context, key string, cbs ...CallbackFunc) int64 {
	result := MustDoCallback(func() (interface{}, error) {
		return QueryInt64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(int64)
}

func QueryUint64(c echo.Context, key string) (uint64, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return 0, NewQueryEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func MustQueryUInt64(c echo.Context, key string, cbs ...CallbackFunc) uint64 {
	result := MustDoCallback(func() (interface{}, error) {
		return QueryUint64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(uint64)
}

func QueryString(c echo.Context, key string) (string, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewQueryEmptyError(key)
	} else {
		return v, nil
	}
}

func MustQueryString(c echo.Context, key string, cbs ...CallbackFunc) string {
	result := MustDoCallback(func() (interface{}, error) {
		return QueryString(c, key)
	}, CodeBadRequest, cbs...)
	return result.(string)
}

func PostFormBool(c echo.Context, key string) (bool, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return false, NewPostFormEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func MustPostFormBool(c echo.Context, key string, cbs ...CallbackFunc) bool {
	result := MustDoCallback(func() (interface{}, error) {
		return PostFormBool(c, key)
	}, CodeBadRequest, cbs...)
	return result.(bool)
}

func PostFormInt64(c echo.Context, key string) (int64, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return 0, NewPostFormEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func MustPostFormInt64(c echo.Context, key string, cbs ...CallbackFunc) int64 {
	result := MustDoCallback(func() (interface{}, error) {
		return PostFormInt64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(int64)
}

func PostFormUint64(c echo.Context, key string) (uint64, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return 0, NewPostFormEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func MustPostFormUInt64(c echo.Context, key string, cbs ...CallbackFunc) uint64 {
	result := MustDoCallback(func() (interface{}, error) {
		return PostFormUint64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(uint64)
}

func PostFormString(c echo.Context, key string) (string, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewPostFormEmptyError(key)
	} else {
		return v, nil
	}
}

func MustPostFormString(c echo.Context, key string, cbs ...CallbackFunc) string {
	result := MustDoCallback(func() (interface{}, error) {
		return PostFormString(c, key)
	}, CodeBadRequest, cbs...)
	return result.(string)
}

func FormBool(c echo.Context, key string) (bool, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return false, NewFormEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func MustFormBool(c echo.Context, key string, cbs ...CallbackFunc) bool {
	result := MustDoCallback(func() (interface{}, error) {
		return FormBool(c, key)
	}, CodeBadRequest, cbs...)
	return result.(bool)
}

func FormInt64(c echo.Context, key string) (int64, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return 0, NewFormEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func MustFormInt64(c echo.Context, key string, cbs ...CallbackFunc) int64 {
	result := MustDoCallback(func() (interface{}, error) {
		return FormInt64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(int64)
}

func FormUint64(c echo.Context, key string) (uint64, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return 0, NewFormEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func MustFormUint64(c echo.Context, key string, cbs ...CallbackFunc) uint64 {
	result := MustDoCallback(func() (interface{}, error) {
		return FormUint64(c, key)
	}, CodeBadRequest, cbs...)
	return result.(uint64)
}

func FormString(c echo.Context, key string) (string, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewFormEmptyError(key)
	} else {
		return v, nil
	}
}

func MustFormString(c echo.Context, key string, cbs ...CallbackFunc) string {
	result := MustDoCallback(func() (interface{}, error) {
		return FormString(c, key)
	}, CodeBadRequest, cbs...)
	return result.(string)
}

// MustFormFile parses multipart message with panic.
// Note: Close needs to be called after MustFormFile is called successfully.
func MustFormFile(c echo.Context, key string, cbs ...CallbackFunc) multipart.File {
	result := MustDoCallback(func() (interface{}, error) {
		if file, err := c.FormFile(key); err != nil {
			return nil, err
		} else {
			return file.Open()
		}
	}, CodeBadRequest, cbs...)
	return result.(multipart.File)
}

func EnvBool(c echo.Context, key string) (bool, error) {
	if v := os.Getenv(key); handy.IsEmptyStr(v) {
		return false, NewEnvEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func MustEnvBool(c echo.Context, key string, cbs ...CallbackFunc) bool {
	result := MustDoCallback(func() (interface{}, error) {
		return EnvBool(c, key)
	}, CodeInternalErr, cbs...)
	return result.(bool)
}

func EnvInt64(c echo.Context, key string) (int64, error) {
	if v := os.Getenv(key); handy.IsEmptyStr(v) {
		return 0, NewEnvEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func MustEnvInt64(c echo.Context, key string, cbs ...CallbackFunc) int64 {
	result := MustDoCallback(func() (interface{}, error) {
		return EnvInt64(c, key)
	}, CodeInternalErr, cbs...)
	return result.(int64)
}

func EnvUint64(c echo.Context, key string) (uint64, error) {
	if v := os.Getenv(key); handy.IsEmptyStr(v) {
		return 0, NewEnvEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func MustEnvUint64(c echo.Context, key string, cbs ...CallbackFunc) uint64 {
	result := MustDoCallback(func() (interface{}, error) {
		return EnvUint64(c, key)
	}, CodeInternalErr, cbs...)
	return result.(uint64)
}

func EnvString(c echo.Context, key string) (string, error) {
	if v := os.Getenv(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewEnvEmptyError(key)
	} else {
		return v, nil
	}
}

func MustEnvString(c echo.Context, key string, cbs ...CallbackFunc) string {
	result := MustDoCallback(func() (interface{}, error) {
		return EnvString(c, key)
	}, CodeInternalErr, cbs...)
	return result.(string)
}

// EncodeValues needs tag "url" in fields of v.
func EncodeValues(v interface{}) (url.Values, error) {
	return query.Values(v)
}

func MustEncodeValues(v interface{}, cbs ...CallbackFunc) url.Values {
	result := MustDoCallback(func() (interface{}, error) {
		return EncodeValues(v)
	}, CodeEncodeErr, cbs...)
	return result.(url.Values)
}

// EncodeJSON needs tag "json" in fields of v.
// Note: ReleaseBuffer needs to be called after EncodeJSON is called successfully.
func EncodeJSON(v interface{}) (*bytes.Buffer, error) {
	buffer := AcquireBuffer()
	if err := json.NewEncoder(buffer).Encode(v); err != nil {
		ReleaseBuffer(buffer)
		return nil, err
	}
	return buffer, nil
}

// MustEncodeJSON needs tag "json" in fields of v.
// Note: ReleaseBuffer needs to be called after MustEncodeJSON is called successfully.
func MustEncodeJSON(v interface{}, cbs ...CallbackFunc) *bytes.Buffer {
	result := MustDoCallback(func() (interface{}, error) {
		return EncodeJSON(v)
	}, CodeEncodeErr, cbs...)
	return result.(*bytes.Buffer)
}

func MustDo(run RunFunc, codes ...int) interface{} {
	code := CodeDownstreamErr
	if len(codes) > 0 {
		code = codes[0]
	}

	return MustDoCallback(run, code)
}

// MustDoCallback will call cbs if and only if run has error.
func MustDoCallback(run RunFunc, code int, cbs ...CallbackFunc) interface{} {
	if run == nil {
		return nil
	}

	result, err := run()
	if err == nil {
		return result
	}

	for _, cb := range cbs {
		cb(err)
	}

	if !IsEchotoolError(err) {
		err = AcquireEchotoolError(code, err)
	}
	panic(err)
}

func GetRequestHost(req *http.Request) (host string) {
	if host = req.Host; handy.IsEmptyStr(host) {
		host = req.URL.Host
	}
	return
}

func GetUUID(c echo.Context) (u string) {
	u, _ = handy.GetUUID()
	return
}

func GetHostname() (host string) {
	host, _ = os.Hostname()
	return
}

func GetHandlerName(f HandlerFunc) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return name[strings.LastIndex(name, handy.StrDot)+1:]
}
