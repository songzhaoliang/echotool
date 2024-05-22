package main

import (
	"fmt"

	vd "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/songzhaoliang/echotool"
	evd "github.com/songzhaoliang/echotool/validator"
)

type User struct {
	ID      int       `json:"id" valid:"userid"`
	Name    string    `json:"name" valid:"required"`
	Gender  string    `json:"gender" valid:"gender|len=0"`
	Courses []*Course `json:"courses" valid:"gte=0,dive"`
}

type Course struct {
	Name  string `json:"name" valid:"required"`
	Score int    `json:"score" valid:"min=0,max=100"`
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
	fs := map[string]vd.Func{
		"userid": IsValidUserID,
		"gender": IsValidGender,
	}

	for k, f := range fs {
		if err := evd.EchotoolValidator.RegisterValidation(k, f); err != nil {
			panic(err)
		}
	}
}

func IsValidUserID(l vd.FieldLevel) bool {
	id := l.Field().Int()

	if id <= 0 {
		return false
	}
	return true
}

func IsValidGender(l vd.FieldLevel) bool {
	gender := l.Field().String()

	switch gender {
	case "male", "female":
		return true
	default:
		return false
	}
}
