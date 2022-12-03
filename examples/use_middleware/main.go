package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := echo.New()
	r.Use(echotool.SetRequestID(echotool.GetUUID))

	e := echotool.NewEngine()
	e.Use(CheckToken)

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CheckToken(c echo.Context, ec *echotool.Context) {
	token := echotool.MustHeaderString(c, "X-Token")
	if token != "RightToken" {
		ec.Abort(echotool.CodeForbidden, fmt.Errorf("invalid token"))
	}
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
