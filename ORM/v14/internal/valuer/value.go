package valuer

import (
	"database/sql"
	"geektime/ORM/v14/model"
)

//因为有 reflect 和 unsafe的实现
//所以用接口

type Value interface {
	SetColumns(rows *sql.Rows) error
	Field(field string) (any, error)
}

type Creator func(val any, model *model.Model) Value
