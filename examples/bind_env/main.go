package main

import (
	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
)

type Credential struct {
	AccessKey string `env:"ACCESS_KEY"`
	Secretkey string `env:"SECRET_KEY"`
}

func main() {
	r := echo.New()
	r.Use(echotool.SetRequestID(echotool.GetUUID))

	e := echotool.NewEngine()

	r.GET("/credentials", e.EchoHandler(GetCredential))

	r.Start(":1323")
}

func GetCredential(c echo.Context, ec *echotool.Context) {
	credential := &Credential{}
	echotool.New(c, credential).BindEnv().MustEnd()

	ec.Finish(echotool.CodeOKZero, credential)
}
