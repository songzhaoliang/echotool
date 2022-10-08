//go:build !reflect2
// +build !reflect2

package binder

import (
	"reflect"
	"strconv"

	"github.com/popeyeio/handy"
)

func Bind(obj interface{}, values map[string][]string, tagKey string, canonical bool) error {
	rt := reflect.TypeOf(obj).Elem()
	rv := reflect.ValueOf(obj).Elem()

	for i := 0; i < rt.NumField(); i++ {
		rtf := rt.Field(i)
		rvf := rv.Field(i)
		if !rvf.CanSet() {
			continue
		}

		kind := rvf.Kind()
		tag := rtf.Tag.Get(tagKey)
		switch tag {
		case handy.StrHyphen:
			continue
		case handy.StrEmpty:
			tag = rtf.Name

			if kind == reflect.Struct {
				if err := Bind(rvf.Addr().Interface(), values, tagKey, canonical); err != nil {
					return err
				}
				continue
			}
		}

		vals, exists := values[canonicalKey(tag, canonical)]
		if !exists {
			continue
		}

		size := len(vals)
		if kind == reflect.Slice && size > 0 {
			elemKind := rvf.Type().Elem().Kind()
			slice := reflect.MakeSlice(rvf.Type(), size, size)
			for j := 0; j < size; j++ {
				if err := SetField(elemKind, vals[j], slice.Index(j)); err != nil {
					return err
				}
			}
			rvf.Set(slice)
		} else if size > 0 {
			if err := SetField(kind, vals[0], rvf); err != nil {
				return err
			}
		}
	}
	return nil
}

func SetField(kind reflect.Kind, val string, field reflect.Value) error {
	switch kind {
	case reflect.Bool:
		return SetBoolField(val, field)
	case reflect.Int:
		return SetIntField(val, 0, field)
	case reflect.Int8:
		return SetIntField(val, 8, field)
	case reflect.Int16:
		return SetIntField(val, 16, field)
	case reflect.Int32:
		return SetIntField(val, 32, field)
	case reflect.Int64:
		return SetIntField(val, 64, field)
	case reflect.Uint:
		return SetUintField(val, 0, field)
	case reflect.Uint8:
		return SetUintField(val, 8, field)
	case reflect.Uint16:
		return SetUintField(val, 16, field)
	case reflect.Uint32:
		return SetUintField(val, 32, field)
	case reflect.Uint64:
		return SetUintField(val, 64, field)
	case reflect.Float32:
		return SetFloatField(val, 32, field)
	case reflect.Float64:
		return SetFloatField(val, 64, field)
	case reflect.String:
		field.SetString(val)
		return nil
	}
	return ErrInvalidType
}

func SetBoolField(val string, field reflect.Value) error {
	v, err := strconv.ParseBool(convertValue(val))
	if err == nil {
		field.SetBool(v)
	}
	return err
}

func SetIntField(val string, bitSize int, field reflect.Value) error {
	v, err := strconv.ParseInt(convertValue(val), 10, bitSize)
	if err == nil {
		field.SetInt(v)
	}
	return err
}

func SetUintField(val string, bitSize int, field reflect.Value) error {
	v, err := strconv.ParseUint(convertValue(val), 10, bitSize)
	if err == nil {
		field.SetUint(v)
	}
	return err
}

func SetFloatField(val string, bitSize int, field reflect.Value) error {
	v, err := strconv.ParseFloat(convertValue(val), bitSize)
	if err == nil {
		field.SetFloat(v)
	}
	return err
}
