package structs

import (
	"reflect"
	"time"

	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"
)

// Option represents whether to target field.
type Option func(srcValue, dstValue reflect.Value, srcField reflect.StructField) bool

func WithIgnoreFields(names ...string) Option {
	return func(_, _ reflect.Value, srcField reflect.StructField) bool {
		return !slices.Contains(names, srcField.Name)
	}
}

func WithTargetFields(names ...string) Option {
	return func(_, _ reflect.Value, srcField reflect.StructField) bool {
		return slices.Contains(names, srcField.Name)
	}
}

func WithValue(v any) Option {
	return func(srcValue, _ reflect.Value, field reflect.StructField) bool {
		// TODO(takashabe): consider field pointer. trying recursive reflect.Indirect while field is pointer.
		return reflect.DeepEqual(v, srcValue.FieldByName(field.Name).Interface())
	}
}

// WithDefaultValue valid when the field in dst is a type-specific default value.
func WithDefaultValue() Option {
	return func(_, dstValue reflect.Value, srcField reflect.StructField) bool {
		return isDefaultValue(dstValue, srcField)
	}
}

// WithIgnoreSourceDefaultValue not valid if the src field is a type-specific default value.
func WithIgnoreSourceDefaultValue() Option {
	return func(srcValue, _ reflect.Value, srcField reflect.StructField) bool {
		return !isDefaultValue(srcValue, srcField)
	}
}

func isDefaultValue(v reflect.Value, field reflect.StructField) bool {
	if v.FieldByName(field.Name) == (reflect.Value{}) {
		return false
	}

	switch field.Type.Kind() {
	case reflect.Bool:
		return !v.FieldByName(field.Name).Bool()
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return v.FieldByName(field.Name).Int() == 0
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return v.FieldByName(field.Name).Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.FieldByName(field.Name).Float() == 0.0
	case reflect.String:
		return v.FieldByName(field.Name).String() == ""
	case reflect.Struct:
		return isDefaultStruct(v.FieldByName(field.Name).Interface())
	case reflect.Pointer:
		return v.FieldByName(field.Name).IsNil()
	case reflect.Slice:
		return v.FieldByName(field.Name).Len() == 0
	default:
		return false
	}
}

// isDefaultStruct returns whether a known structure is the default value.
func isDefaultStruct(a any) bool {
	switch a.(type) {
	case time.Time:
		return a == time.Time{}
	case null.Bool:
		return a == null.Bool{}
	case null.Int:
		return a == null.Int{}
	case null.Float:
		return a == null.Float{}
	case null.Time:
		return a == null.Time{}
	default:
		return false
	}
}
