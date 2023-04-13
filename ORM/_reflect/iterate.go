package _reflect

import (
	"errors"
	"reflect"
)

// Iterate 遍历 array slice string
func Iterate(entity any) ([]any, error) {
	val := reflect.ValueOf(entity)
	kind := val.Kind()
	if kind != reflect.Array && kind != reflect.Slice && kind != reflect.String {
		return nil, errors.New("非法类型")
	}
	res := make([]any, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		res = append(res, elem.Interface())
	}
	return res, nil
}

// IterateMapV1 返回键，值
func IterateMapV1(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	kind := val.Kind()
	if kind != reflect.Map {
		return nil, nil, errors.New("非法类型")
	}
	l := val.Len()
	keys := make([]any, 0, l)
	values := make([]any, 0, l)
	for _, k := range val.MapKeys() {
		//k是key
		keys = append(keys, k.Interface())
		//再通过key找到value
		v := val.MapIndex(k)
		values = append(values, v.Interface())
	}
	return keys, values, nil
}

// IterateMapV2 返回键，值
func IterateMapV2(entity any) ([]any, []any, error) {
	val := reflect.ValueOf(entity)
	kind := val.Kind()
	if kind != reflect.Map {
		return nil, nil, errors.New("非法类型")
	}
	l := val.Len()
	keys := make([]any, 0, l)
	values := make([]any, 0, l)
	itr := val.MapRange()
	for itr.Next() {
		keys = append(keys, itr.Key().Interface())
		values = append(values, itr.Value().Interface())
	}
	return keys, values, nil
}
