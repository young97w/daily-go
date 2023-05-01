package valuer

import (
	"database/sql"
	"geektime/ORM/v10/internal/model"
)

//因为有 reflect 和 unsafe的实现
//所以用接口

type Value interface {
	SetColumns(rows *sql.Rows) error
}

type Creator func(val any, model *model.Model) Value
