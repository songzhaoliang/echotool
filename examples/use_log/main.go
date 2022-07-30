package main

import (
	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	defer echotool.FlushLog()

	r := echo.New()
	r.Use(echotool.SetRequestID(echotool.GetUUID))

	e := echotool.NewEngine()
	e.Use(echotool.AddTraceID(echotool.GetRequestID))
	e.Use(echotool.AddNotice("host", echotool.GetHostname()))

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	echotool.CtxInfo(ec, "user is %+v", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
