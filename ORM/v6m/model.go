package v1

import (
	"errors"
	"geektime/ORM/internal/errs"
	"reflect"
)

// 我们支持的全部标签上的 key 都放在这里
// 方便用户查找，和我们后期维护
const (
	tagKeyColumn = "column"
)

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}

type Model struct {
	tableName string
	fields    map[string]*field
	columns   map[string]*field
}

type field struct {
	fieldName string
	colName   string
	typ       reflect.Type
}

type ModelOpt func(model *Model) error

func ModelWithTableName(tableName string) ModelOpt {
	return func(model *Model) error {
		if tableName == "" {
			return errors.New("orm: table name can't be empty")
		}
		model.tableName = tableName
		return nil
	}
}

func ModelWithColumnName(field, colName string) ModelOpt {
	return func(model *Model) error {
		f, ok := model.fields[field]
		if !ok || colName == "" {
			return errs.NewErrUnknownField(field)
		}
		f.colName = colName
		return nil
	}
}
