package echotool

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommonResponse struct {
	RequestID string      `json:"request_id,omitempty"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

func RespOK(id string, code int, data interface{}) *CommonResponse {
	return &CommonResponse{
		RequestID: id,
		Code:      code,
		Message:   CodeMsg(code),
		Data:      data,
	}
}

func RespError(id string, code int, err error) *CommonResponse {
	resp := &CommonResponse{
		RequestID: id,
		Code:      code,
		Message:   CodeMsg(code),
	}
	if err != nil {
		resp.Message += " - " + err.Error()
	}
	return resp
}

func FinishWithCodeData(c echo.Context, code int, data interface{}) {
	status := HTTPStatus(code)
	if status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect {
		c.JSON(status, RespOK(GetRequestID(c), code, data))
	}
	c.Redirect(status, data.(string))
}

func AbortWithCodeErr(c echo.Context, code int, err error) {
	c.JSON(HTTPStatus(code), RespError(GetRequestID(c), code, err))
}

func GetCommonFinisher() HandlerFunc {
	return func(c echo.Context, ec *Context) {
		FinishWithCodeData(c, ec.GetCode(), ec.GetData())
	}
}

func GetCommonAborter() HandlerFunc {
	return func(c echo.Context, ec *Context) {
		AbortWithCodeErr(c, ec.GetCode(), ec.GetError())
	}
}
