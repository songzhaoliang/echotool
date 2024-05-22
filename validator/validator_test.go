package validator

import (
	"testing"

	vd "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID   int    `valid:"userid"`
	Name string `valid:"username"`
}

func TestEchotoolValidator(t *testing.T) {
	fs := map[string]vd.Func{
		"userid":   IsValidUserID,
		"username": IsValidUsername,
	}

	for k, f := range fs {
		if err := EchotoolValidator.RegisterValidation(k, f); err != nil {
			assert.NoError(t, err)
		}
	}

	err := EchotoolValidator.ValidateStruct(&User{1, "peter"})
	assert.NoError(t, err)
}

func IsValidUserID(l vd.FieldLevel) bool {
	id := l.Field().Int()

	if id <= 0 {
		return false
	}
	return true
}

func IsValidUsername(l vd.FieldLevel) bool {
	name := l.Field().String()

	length := len(name)
	if length < 2 || length > 50 {
		return false
	}
	return true
}
