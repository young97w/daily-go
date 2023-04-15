package v1

import (
	"geektime/ORM/internal/errs"
	"reflect"
	"unicode"
)

func (r *registry) parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	numField := typ.NumField()
	fields := make(map[string]*field, numField)

	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		fields[f.Name] = &field{colName: underscoreName(f.Name)}
	}

	return &model{
		tableName: underscoreName(typ.Name()),
		fields:    fields,
	}, nil
}

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

func underscoreName(s string) string {
	var buf []byte
	for i, v := range s {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
