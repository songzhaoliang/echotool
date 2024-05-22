package validator

import (
	"reflect"
	"sync"

	vd "github.com/go-playground/validator/v10"
)

var EchotoolValidator = &echotoolValidator{}

type echotoolValidator struct {
	once     sync.Once
	validate *vd.Validate
}

var _ Validator = (*echotoolValidator)(nil)

func (v *echotoolValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyInit()
		return v.validate.Struct(obj)
	}
	return nil
}

func (v *echotoolValidator) RegisterValidation(key string, fn vd.Func) error {
	v.lazyInit()
	return v.validate.RegisterValidation(key, fn)
}

func (v *echotoolValidator) RegisterAlias(alias, tags string) {
	v.lazyInit()
	v.validate.RegisterAlias(alias, tags)
}

func (v *echotoolValidator) RegisterStructValidation(fn vd.StructLevelFunc, types ...interface{}) {
	v.lazyInit()
	v.validate.RegisterStructValidation(fn, types...)
}

func (v *echotoolValidator) RegisterCustomTypeFunc(fn vd.CustomTypeFunc, types ...interface{}) {
	v.lazyInit()
	v.validate.RegisterCustomTypeFunc(fn, types...)
}

func (v *echotoolValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = vd.New()
		v.validate.SetTagName(TagValid)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	rv := reflect.ValueOf(data)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		kind = rv.Elem().Kind()
	}
	return kind
}
