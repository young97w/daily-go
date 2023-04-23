package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type FieldMeta struct {
	Offset uintptr
	typ    reflect.Type
}

type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numField := typ.NumField()
	fields := make(map[string]FieldMeta, numField)

	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		fields[f.Name] = FieldMeta{
			Offset: f.Offset,
			typ:    f.Type,
		}
	}

	return &UnsafeAccessor{
		fields:  fields,
		address: reflect.ValueOf(entity).UnsafePointer(),
	}
}

func (a *UnsafeAccessor) field(field string) (any, error) {
	f, ok := a.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}

	fAddress := unsafe.Pointer(uintptr(a.address) + f.Offset)
	return reflect.NewAt(f.typ, fAddress).Elem().Interface(), nil
}

func (a *UnsafeAccessor) setField(field string, val any) error {
	f, ok := a.fields[field]
	if !ok {
		return errors.New("非法字段")
	}

	fAddress := unsafe.Pointer(uintptr(a.address) + f.Offset)
	reflect.NewAt(f.typ, fAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}
