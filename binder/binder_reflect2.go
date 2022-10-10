//go:build reflect2
// +build reflect2

package binder

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/modern-go/concurrent"
	"github.com/modern-go/reflect2"
	"github.com/popeyeio/handy"
)

func Bind(obj interface{}, values map[string][]string, tagKey string, canonical bool) error {
	rt := getRTypeFromCache(obj)
	ptr := reflect2.PtrOf(obj)

	for i := 0; i < rt.NumField(); i++ {
		rtf := rt.Field(i)
		fptr := rtf.UnsafeGet(ptr)
		typ := rtf.Type()
		kind := typ.Kind()
		tag := rtf.Tag().Get(tagKey)

		switch tag {
		case handy.StrHyphen:
			continue
		case handy.StrEmpty:
			tag = rtf.Name()

			if kind == reflect.Struct {
				if err := Bind(typ.PackEFace(fptr), values, tagKey, canonical); err != nil {
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
			sliceType := typ.(*reflect2.UnsafeSliceType)
			elemKind := sliceType.Elem().Kind()
			sliceType.UnsafeSet(fptr, sliceType.UnsafeMakeSlice(size, size))
			for j := 0; j < size; j++ {
				if err := SetField(elemKind, vals[j], sliceType.UnsafeGetIndex(fptr, j)); err != nil {
					return err
				}
			}
		} else if size > 0 {
			if err := SetField(kind, vals[0], fptr); err != nil {
				return err
			}
		}
	}
	return nil
}

func SetField(kind reflect.Kind, val string, ptr unsafe.Pointer) error {
	switch kind {
	case reflect.Bool:
		return SetBoolField(val, ptr)
	case reflect.Int:
		return SetIntField(val, ptr)
	case reflect.Int8:
		return SetInt8Field(val, ptr)
	case reflect.Int16:
		return SetInt16Field(val, ptr)
	case reflect.Int32:
		return SetInt32Field(val, ptr)
	case reflect.Int64:
		return SetInt64Field(val, ptr)
	case reflect.Uint:
		return SetUintField(val, ptr)
	case reflect.Uint8:
		return SetUint8Field(val, ptr)
	case reflect.Uint16:
		return SetUint16Field(val, ptr)
	case reflect.Uint32:
		return SetUint32Field(val, ptr)
	case reflect.Uint64:
		return SetUint64Field(val, ptr)
	case reflect.Float32:
		return SetFloat32Field(val, ptr)
	case reflect.Float64:
		return SetFloat64Field(val, ptr)
	case reflect.String:
		SetString(val, ptr)
		return nil
	}
	return ErrInvalidType
}

func SetBoolField(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseBool(convertValue(val))
	if err == nil {
		*(*bool)(ptr) = v
	}
	return err
}

func SetIntField(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseInt(convertValue(val), 10, 0)
	if err == nil {
		*(*int)(ptr) = int(v)
	}
	return err
}

func SetInt8Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseInt(convertValue(val), 10, 8)
	if err == nil {
		*(*int8)(ptr) = int8(v)
	}
	return err
}

func SetInt16Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseInt(convertValue(val), 10, 16)
	if err == nil {
		*(*int16)(ptr) = int16(v)
	}
	return err
}

func SetInt32Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseInt(convertValue(val), 10, 32)
	if err == nil {
		*(*int32)(ptr) = int32(v)
	}
	return err
}

func SetInt64Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseInt(convertValue(val), 10, 64)
	if err == nil {
		*(*int64)(ptr) = v
	}
	return err
}

func SetUintField(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseUint(convertValue(val), 10, 0)
	if err == nil {
		*(*uint)(ptr) = uint(v)
	}
	return err
}

func SetUint8Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseUint(convertValue(val), 10, 8)
	if err == nil {
		*(*uint8)(ptr) = uint8(v)
	}
	return err
}

func SetUint16Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseUint(convertValue(val), 10, 16)
	if err == nil {
		*(*uint16)(ptr) = uint16(v)
	}
	return err
}

func SetUint32Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseUint(convertValue(val), 10, 32)
	if err == nil {
		*(*uint32)(ptr) = uint32(v)
	}
	return err
}

func SetUint64Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseUint(convertValue(val), 10, 64)
	if err == nil {
		*(*uint64)(ptr) = v
	}
	return err
}

func SetFloat32Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseFloat(convertValue(val), 32)
	if err == nil {
		*(*float32)(ptr) = float32(v)
	}
	return err
}

func SetFloat64Field(val string, ptr unsafe.Pointer) error {
	v, err := strconv.ParseFloat(convertValue(val), 64)
	if err == nil {
		*(*float64)(ptr) = v
	}
	return err
}

func SetString(val string, ptr unsafe.Pointer) {
	*(*string)(ptr) = val
}

var cache = concurrent.NewMap()

func getRTypeFromCache(obj interface{}) (rt *reflect2.UnsafeStructType) {
	key := reflect2.RTypeOf(obj)
	if val, exists := cache.Load(key); exists {
		rt = val.(*reflect2.UnsafeStructType)
	} else {
		rt = reflect2.Type2(reflect.TypeOf(obj).Elem()).(*reflect2.UnsafeStructType)
		cache.Store(key, rt)
	}
	return
}
