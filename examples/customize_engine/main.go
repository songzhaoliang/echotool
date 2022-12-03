package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool"
)

type Response struct {
	CodeN int         `json:"CodeN"`
	Code  string      `json:"Code"`
	Data  interface{} `json:"Data,omitempty"`
}

func RespOK(code int, data interface{}) *Response {
	return &Response{
		CodeN: code,
		Code:  echotool.CodeMsg(code),
		Data:  data,
	}
}

func RespError(code int, err error) *Response {
	return &Response{
		CodeN: code,
		Code:  fmt.Sprintf("%s - %v", echotool.CodeMsg(code), err),
	}
}

func NewEngine() *echotool.Engine {
	opts := []echotool.Option{
		echotool.WithFinisher(func(c echo.Context, ec *echotool.Context) {
			code, data := ec.GetCode(), ec.GetData()
			c.JSON(echotool.HTTPStatus(code), RespOK(code, data))
		}),
		echotool.WithAborter(func(c echo.Context, ec *echotool.Context) {
			code, err := ec.GetCode(), ec.GetError()
			fmt.Printf("%s - %v\n", echotool.CodeMsg(code), err)
			c.JSON(echotool.HTTPStatus(code), RespError(code, err))
		}),
	}
	return echotool.NewEngine(opts...)
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := echo.New()

	e := NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
