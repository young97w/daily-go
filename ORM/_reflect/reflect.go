package _reflect

import (
	"errors"
	"reflect"
)

func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持nil")
	}
	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)
	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}

	//解开指针
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持类型")
	}

	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		//get type of field
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()

		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}
	return res, nil
}

func SetValue(entity any, field string, newValue any) error {
	val := reflect.ValueOf(entity)

	for val.Type().Kind() == reflect.Pointer {
		val = val.Elem()
	}
	fieldVal := val.FieldByName(field)
	if !fieldVal.CanSet() {
		return errors.New("不可更改字段")
	}
	fieldVal.Set(reflect.ValueOf(newValue))
	return nil
}
