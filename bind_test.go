package echotool

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	vd "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	evd "github.com/songzhaoliang/echotool/validator"
	"github.com/stretchr/testify/assert"
)

const defaultUserID = 100

type User struct {
	ID       int    `header:"X-Id" valid:"userid"`
	Surname  string `form:"surname"`
	Name     string `form:"name" valid:"username"`
	FullName string
}

func (u *User) BeforeBind(c echo.Context) error {
	u.ID = defaultUserID
	return nil
}

func (u *User) AfterBind(c echo.Context) error {
	u.FullName = u.Surname + " " + u.Name
	return nil
}

func TestBind_FirstMethod(t *testing.T) {
	c, u := Preprocess(t)

	err := New(c, u).BindHeader().FormBindBody().Validate().End()

	Postprocess(t, u, err)
}

func TestBind_SecondMethod(t *testing.T) {
	c, u := Preprocess(t)

	err := Bind(c, u, BHeader|BFormBody|BValidator)

	Postprocess(t, u, err)
}

func Preprocess(t *testing.T) (echo.Context, *User) {
	fs := map[string]vd.Func{
		"userid":   IsValidUserID,
		"username": IsValidUserName,
	}

	for k, f := range fs {
		if err := evd.EchotoolValidator.RegisterValidation(k, f); err != nil {
			assert.NoError(t, err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("surname=li&name=si"))
	req.Header.Set("X-Id", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	return c, &User{}
}

func Postprocess(t *testing.T, u *User, err error) {
	assert.NoError(t, err)
	assert.Equal(t, 1, u.ID)
	assert.Equal(t, "li", u.Surname)
	assert.Equal(t, "si", u.Name)
	assert.Equal(t, "li si", u.FullName)
}

func IsValidUserID(l vd.FieldLevel) bool {
	id := l.Field().Int()

	if id <= 0 {
		return false
	}
	return true
}

func IsValidUserName(l vd.FieldLevel) bool {
	name := l.Field().String()

	length := len(name)
	if length < 2 || length > 50 {
		return false
	}
	return true
}
