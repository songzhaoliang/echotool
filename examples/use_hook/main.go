package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool"
	"github.com/songzhaoliang/echotool/binder"
)

const (
	defaultUserID = 100
)

type User struct {
	ID       int    `json:"id"`
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	FullName string `json:"-"`
}

var _ binder.BeforeBinder = (*User)(nil)
var _ binder.AfterBinder = (*User)(nil)

func (u *User) BeforeBind(c echo.Context) error {
	u.ID = defaultUserID
	return nil
}

func (u *User) AfterBind(c echo.Context) error {
	u.FullName = u.Surname + " " + u.Name
	return nil
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
