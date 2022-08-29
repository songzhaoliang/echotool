package validator

import (
	rf "reflect"
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

func IsValidUserID(_ *vd.Validate, _, _, v rf.Value, _ rf.Type, _ rf.Kind, _ string) bool {
	id := v.Int()

	if id <= 0 {
		return false
	}
	return true
}

func IsValidUserName(_ *vd.Validate, _, _, v rf.Value, _ rf.Type, _ rf.Kind, _ string) bool {
	name := v.String()

	length := len(name)
	if length < 2 || length > 50 {
		return false
	}
	return true
}
