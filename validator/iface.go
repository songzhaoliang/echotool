package validator

import (
	vd "github.com/go-playground/validator/v10"
)

const (
	TagValid = "valid"
)

type Validator interface {
	ValidateStruct(obj interface{}) error
	RegisterValidation(key string, fn vd.Func) error
	RegisterAlias(alias, tags string)
	RegisterStructValidation(fn vd.StructLevelFunc, types ...interface{})
	RegisterCustomTypeFunc(fn vd.CustomTypeFunc, types ...interface{})
}
