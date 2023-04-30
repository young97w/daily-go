package model

import (
	"errors"
	"geektime/ORM/internal/errs"
	"reflect"
)

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	TagKeyColumn = "column"
)

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}

type Model struct {
	TableName string
	Fields    map[string]*Field
	Columns   map[string]*Field
	ColSlice  []*Field
}

type Field struct {
	FieldName string
	ColName   string
	Typ       reflect.Type
	//offset 在结构体中的相对地址偏移量
	Offset uintptr
}

type ModelOpt func(model *Model) error

func ModelWithTableName(tableName string) ModelOpt {
	return func(model *Model) error {
		if tableName == "" {
			return errors.New("orm: table name can't be empty")
		}
		model.TableName = tableName
		return nil
	}
}

func ModelWithColumnName(field, colName string) ModelOpt {
	return func(model *Model) error {
		f, ok := model.Fields[field]
		if !ok || colName == "" {
			return errs.NewErrUnknownField(field)
		}
		f.ColName = colName
		return nil
	}
}
