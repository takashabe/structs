package structs

import "reflect"

var DiffTagName = "json"

// DiffFields returns different field names
func DiffFields[T, T2 any](src T, dst T2, opts ...Option) ([]string, error) {
	sv, err := getStruct(reflect.ValueOf(src))
	if err != nil {
		return nil, err
	}
	dv, err := getStruct(reflect.ValueOf(dst))
	if err != nil {
		return nil, err
	}

	diff := []string{}
	fields := reflect.VisibleFields(sv.Type())
	for _, f := range fields {
		if !f.IsExported() {
			continue
		}

		canSet := true
		for _, opt := range opts {
			ok := opt(sv, dv, f)
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
		diffField := f.Name
		if tag, ok := f.Tag.Lookup(DiffTagName); ok {
			diffField = tag
		}
		diff = append(diff, diffField)
	}
	return diff, nil
}
