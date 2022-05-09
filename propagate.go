package structs

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"
)

// PropagateOption represents whether to copy field.
type PropagateOption func(value reflect.Value, field reflect.StructField) bool

func PropagateWithIgnoreFields(names ...string) PropagateOption {
	return func(value reflect.Value, field reflect.StructField) bool {
		return !slices.Contains(names, field.Name)
	}
}

func PropagateWithValue(v any) PropagateOption {
	return func(value reflect.Value, field reflect.StructField) bool {
		// TODO(takashabe): consider field pointer. trying recursive reflect.Indirect while field is pointer.
		return reflect.DeepEqual(v, value.FieldByName(field.Name).Interface())
	}
}

// PropagateValues copy from `src` field values to `dst`.
func PropagateValues[T, T2 any](src T, dst T2, opts ...PropagateOption) (T2, error) {
	sv, err := getStruct(reflect.ValueOf(src))
	if err != nil {
		return dst, err
	}
	dv, err := getStruct(reflect.ValueOf(dst))
	if err != nil {
		return dst, err
	}

	fields := reflect.VisibleFields(sv.Type())
	for _, f := range fields {
		if !f.IsExported() {
			continue
		}

		canSet := true
		for _, opt := range opts {
			ok := opt(sv, f)
			if !ok {
				canSet = false
				break
			}
		}
		if !canSet {
			continue
		}

		df := dv.FieldByName(f.Name)
		// field not exists
		if df == (reflect.Value{}) {
			continue
		}
		df.Set(sv.FieldByName(f.Name))
	}
	return dst, nil
}

func getStruct(rv reflect.Value) (reflect.Value, error) {
	switch rv.Kind() {
	case reflect.Pointer:
		e := rv.Elem()
		if e.Kind() == reflect.Struct {
			return e, nil
		}
		fallthrough
	default:
		return rv, fmt.Errorf("value want struct pointer, got %s", rv.Kind().String())
	}
}
