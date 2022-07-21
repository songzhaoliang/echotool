package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
)

const (
	CodeOracleErr = 50018
)

var codeMsg = map[int]string{
	CodeOracleErr: "oracle error",
}

func init() {
	echotool.RegisterCode(CodeOracleErr, codeMsg[CodeOracleErr], http.StatusInternalServerError)
}

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

	MustOracleCreateUser(user)

	ec.Finish(echotool.CodeOKZero, nil)
}

func OracleCreateUser(user *User) error {
	return fmt.Errorf("duplicate entry")
}

func MustOracleCreateUser(user *User) {
	echotool.MustDo(func() (data interface{}, err error) {
		if err = OracleCreateUser(user); err != nil {
			fmt.Printf("OracleCreateUser error - %v\n", err)
		}
		return
	}, CodeOracleErr)
}
