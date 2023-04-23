package valuer

import (
	"database/sql"
	"geektime/ORM/internal/errs"
	"geektime/ORM/v7/internal/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	addr  unsafe.Pointer
	model *model.Model
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(val any, model *model.Model) Value {
	return unsafeValue{
		addr:  reflect.ValueOf(val).UnsafePointer(),
		model: model,
	}
}

func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	//准备好val []any 来放置每一个 field
	//使用scan 来设置值

	//先校验
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(u.model.Columns) != len(cs) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]any, len(cs))
	//给每一个val 一个指针
	for i, c := range cs {
		f, ok := u.model.Columns[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		ptr := unsafe.Pointer(uintptr(u.addr) + f.Offset)
		colValues[i] = reflect.NewAt(f.Typ, ptr).Interface()
	}

	//使用scan
	return rows.Scan(colValues...)

}
