package validator

import (
	"gopkg.in/go-playground/validator.v8"
)

const (
	TagValid = "valid"
)

type Validator interface {
	ValidateStruct(obj interface{}) error
	RegisterValidation(key string, fn validator.Func) error
	RegisterAliasValidation(alias, tags string)
	RegisterStructValidation(fn validator.StructLevelFunc, types ...interface{})
	RegisterCustomTypeFunc(fn validator.CustomTypeFunc, types ...interface{})
}
