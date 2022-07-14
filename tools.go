package echotool

import (
	"bytes"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
	"github.com/labstack/echo"
	"github.com/popeyeio/handy"
	"github.com/songzhaoliang/echotool/json"
)

func HeaderBool(c echo.Context, key string) (bool, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return false, NewHeaderEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func HeaderInt64(c echo.Context, key string) (int64, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return 0, NewHeaderEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func HeaderUint64(c echo.Context, key string) (uint64, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return 0, NewHeaderEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func HeaderString(c echo.Context, key string) (string, error) {
	if v := c.Request().Header.Get(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewHeaderEmptyError(key)
	} else {
		return v, nil
	}
}

func ParamBool(c echo.Context, key string) (bool, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return false, NewParamEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func ParamInt64(c echo.Context, key string) (int64, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return 0, NewParamEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func ParamUint64(c echo.Context, key string) (uint64, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return 0, NewParamEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func ParamString(c echo.Context, key string) (string, error) {
	if v := c.Param(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewParamEmptyError(key)
	} else {
		return v, nil
	}
}

func QueryBool(c echo.Context, key string) (bool, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return false, NewQueryEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func QueryInt64(c echo.Context, key string) (int64, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return 0, NewQueryEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func QueryUint64(c echo.Context, key string) (uint64, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return 0, NewQueryEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func QueryString(c echo.Context, key string) (string, error) {
	if v := c.QueryParam(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewQueryEmptyError(key)
	} else {
		return v, nil
	}
}

func PostFormBool(c echo.Context, key string) (bool, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return false, NewPostFormEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func PostFormInt64(c echo.Context, key string) (int64, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return 0, NewPostFormEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func PostFormUint64(c echo.Context, key string) (uint64, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return 0, NewPostFormEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func PostFormString(c echo.Context, key string) (string, error) {
	if v := c.Request().PostFormValue(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewPostFormEmptyError(key)
	} else {
		return v, nil
	}
}

func FormBool(c echo.Context, key string) (bool, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return false, NewFormEmptyError(key)
	} else {
		return strconv.ParseBool(v)
	}
}

func FormInt64(c echo.Context, key string) (int64, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return 0, NewFormEmptyError(key)
	} else {
		return strconv.ParseInt(v, 10, 64)
	}
}

func FormUint64(c echo.Context, key string) (uint64, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return 0, NewFormEmptyError(key)
	} else {
		return strconv.ParseUint(v, 10, 64)
	}
}

func FormString(c echo.Context, key string) (string, error) {
	if v := c.FormValue(key); handy.IsEmptyStr(v) {
		return handy.StrEmpty, NewFormEmptyError(key)
	} else {
		return v, nil
	}
}

// EncodeValues needs tag "url" in fields of v.
func EncodeValues(v interface{}) (url.Values, error) {
	return query.Values(v)
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
