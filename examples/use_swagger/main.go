package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool"
	_ "github.com/songzhaoliang/echotool/examples/use_swagger/static"
	"github.com/songzhaoliang/echotool/swagger"
)

/*
	1. download swagger sdk and install swagger command
	   >> go get -u github.com/swaggo/swag
	   >> go install github.com/swaggo/swag/cmd/swag@latest

	2. add api annotations

	3. run swagger command in root folder which contains main.go
	   >> swag init --output static (default: docs)

	4. import swagger output directory
	   >> import _ "github.com/songzhaoliang/echotool/examples/use_swagger/static"

	5. regester swagger router
	   >> swagger.Register(r)
*/

//	@title			Swagger Example
//	@version		1.0
//	@description	This is a server for swagger example
//	@termsOfService	xxx

//	@contact.name	songzhaoliang
//	@contact.url	https://github.com/songzhaoliang/echotool
//	@contact.email	957687172@qq.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:1323
//	@BasePath	/
func main() {
	r := echo.New()
	swagger.Register(r)

	e := echotool.NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CommonResponse struct {
	RequestID string      `json:"request_id,omitempty"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

// CreateUser godoc
//	@Tags			user
//	@Router			/users [post]
//	@Summary		create a user
//	@Description	create a user by id and name
//	@Accept			json
//	@Param			_	body	User	true	"user information"
//	@Produce		json
//	@Success		200	{object}	CommonResponse
//	@Failure		500	{object}	CommonResponse
func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).JSONBindBody().MustEnd()

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}
