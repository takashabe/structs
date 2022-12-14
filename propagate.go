package structs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
)

// PropagateValues copy from `src` field values to `dst`.
func PropagateValues[T, T2 any](src T, dst T2, opts ...Option) (T2, error) {
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
		isTarget := true
		for _, opt := range opts {
			ok := opt(sv, dv, f)
			if !ok {
				isTarget = false
				break
			}
		}
		if !isTarget {
			continue
		}

		if err := set(sv.FieldByName(f.Name), dv.FieldByName(f.Name)); err != nil {
			return dst, fmt.Errorf("failed to set `%s` field. from type `%s` to `%s`: %w", f.Name, f.Type.String(), dv.FieldByName(f.Name).Type().String(), err)
		}
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

// set sets `sv` value to `dv` value.
//
// return value is whether they failed where they should have succeeded.
func set(sv, dv reflect.Value) (err error) {
	defer func() {
		p := recover()
		if p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()

	if !sv.IsValid() {
		return nil
	}

	// field not exists
	if dv == (reflect.Value{}) {
		return nil
	}

	if sv.Kind() == reflect.Pointer && dv.Kind() == reflect.Pointer {
		if sv.IsNil() && dv.IsNil() {
			return nil
		}

		// set nil
		if sv.IsNil() {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}
	}

	// depointer
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if dv.Kind() == reflect.Ptr {
		if dv.IsNil() {
			// allocate default value
			dv.Set(reflect.New(dv.Type().Elem()))
		}
		dv = dv.Elem()
	}

	// assignable primitive type
	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Struct && sv.Kind() != reflect.Map {
		if sv.Type().AssignableTo(dv.Type()) || sv.Type().ConvertibleTo(dv.Type()) {
			if sv.Kind() != reflect.Pointer {
				dv.Set(sv.Convert(dv.Type()))
			} else {
				swap := reflect.New(sv.Type())
				swap.Set(sv.Elem())
				dv.Set(swap.Convert(dv.Type()))
			}
			return nil
		}
	}

	if dv.Kind() == reflect.Interface && reflect.TypeOf(sv.Interface()) != nil {
		dv.Set(reflect.New(reflect.TypeOf(sv.Interface())).Elem())
		return nil
	}

	if sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	if sv.Kind() == reflect.Slice || dv.Kind() == reflect.Slice {
		st := dv.Type()
		v1 := reflect.New(st)
		v1.Elem().Set(reflect.MakeSlice(st, sv.Len(), sv.Cap()))
		dv.Set(v1.Elem())
		itr := v1.Interface()

		if err := copier.CopyWithOption(itr, sv.Interface(), copier.Option{
			IgnoreEmpty: true,
			DeepCopy:    true,
		}); err != nil {
			return fmt.Errorf("failed copier.Copy: %w", err)
		}
		return nil
	}

	if sv.Kind() == reflect.Struct && dv.Kind() == reflect.Struct {
		// allocate default value
		dv.Set(reflect.New(reflect.TypeOf(dv.Interface())).Elem())

		fields := reflect.VisibleFields(sv.Type())
		for _, f := range fields {
			if !f.IsExported() {
				continue
			}
			if err := set(sv.FieldByName(f.Name), dv.FieldByName(f.Name)); err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New("unavailable values")
}
