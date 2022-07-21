package main

import (
	"fmt"

	"github.com/labstack/echo"
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

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()
	// echotool.MustBind(c, user, echotool.BJSONBody)

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
