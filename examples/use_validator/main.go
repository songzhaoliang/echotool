package main

import (
	"fmt"
	"reflect"

	"github.com/labstack/echo"
	"github.com/songzhaoliang/echotool"
	evd "github.com/songzhaoliang/echotool/validator"
	vd "gopkg.in/go-playground/validator.v8"
)

type User struct {
	ID   int    `json:"id" valid:"userid"`
	Name string `json:"name" valid:"required"`
}

func main() {
	InitValidators()

	r := echo.New()
	r.Use(echotool.SetRequestID(echotool.GetUUID))

	e := echotool.NewEngine()

	r.POST("/users", e.EchoHandler(CreateUser))

	r.Start(":1323")
}

func CreateUser(c echo.Context, ec *echotool.Context) {
	user := &User{}
	echotool.New(c, user).JSONBindBody().Validate().MustEnd()
	// echotool.MustBind(c, user, echotool.BJSONBody|echotool.BValidator)

	fmt.Printf("user is %+v\n", user)

	ec.Finish(echotool.CodeOKZero, nil)
}

func InitValidators() {
	validatorFuncs := map[string]vd.Func{
		"userid": IsValidUserID,
	}

	for k, f := range validatorFuncs {
		if err := evd.EchotoolValidator.RegisterValidation(k, f); err != nil {
			panic(err)
		}
	}
}

func IsValidUserID(v *vd.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	id := field.Int()

	if id <= 0 {
		return false
	}
	return true
}
