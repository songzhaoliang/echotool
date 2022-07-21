package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
)

type User struct {
	ID   int    `param:"id"`
	Name string `param:"name"`
}

func main() {
	r := echo.New()
	r.Use(echotool.SetRequestID(echotool.GetUUID))

	e := echotool.NewEngine()

	r.POST("/users/:id/:name", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).BindParam().MustEnd()
	// echotool.MustBind(c, user, echotool.BParam)

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
