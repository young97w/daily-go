package valuer

import (
	"database/sql"
	v1 "geektime/ORM/v7/internal/model"
)

//因为有 reflect 和 unsafe的实现
//所以用接口

type Value interface {
	SetColumns(rows *sql.Rows) error
}

type Creator func(val any, model *v1.Model) Value
