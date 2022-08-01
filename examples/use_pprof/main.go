package main

import (
	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool/pprof"
)

func main() {
	r := echo.New()
	pprof.Register(r)

	v1 := r.Group("/echotool/v1")
	pprof.RouterRegister(v1)

	r.Start(":1323")
}
