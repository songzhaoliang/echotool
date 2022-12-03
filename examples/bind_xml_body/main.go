package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool"
)

type User struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

func main() {
	r := echo.New()
	r.Use(echotool.SetRequestID(echotool.GetUUID))

	e := echotool.NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).XMLBindBody().MustEnd()
	// echotool.MustBind(c, user, echotool.BXMLBody)

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
