package validator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	vd "gopkg.in/go-playground/validator.v8"
)

type User struct {
	ID   int    `valid:"userid"`
	Name string `valid:"username"`
}

func TestEchotoolValidator(t *testing.T) {
	validatorFuncs := map[string]vd.Func{
		"userid":   IsValidUserID,
		"username": IsValidUserName,
	}

	for k, f := range validatorFuncs {
		if err := EchotoolValidator.RegisterValidation(k, f); err != nil {
			assert.NoError(t, err)
		}
	}

	err := EchotoolValidator.ValidateStruct(&User{1, "peter"})
	assert.NoError(t, err)
}

func IsValidUserID(v *vd.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	id := field.Int()

	if id <= 0 {
		return false
	}
	return true
}

func IsValidUserName(v *vd.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	name := field.String()

	length := len(name)
	if length < 2 || length > 50 {
		return false
	}
	return true
}
