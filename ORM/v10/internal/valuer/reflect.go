package valuer

import (
	"database/sql"
	"geektime/ORM/internal/errs"
	"geektime/ORM/v10/internal/model"
	"reflect"
)

type reflectValue struct {
	val   reflect.Value
	model *model.Model
}

var _ Creator = NewReflectValue

func NewReflectValue(val any, model *model.Model) Value {
	return reflectValue{
		val:   reflect.ValueOf(val).Elem(),
		model: model,
	}
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	//使用反射获取类型
	//需要建两个slice
	//scan到一个slice，然后从slice中再赋值到val
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(r.model.Columns) != len(cs) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]any, len(cs))
	colEleValues := make([]reflect.Value, len(cs))
	for i, c := range cs {
		f, ok := r.model.Columns[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		val := reflect.New(f.Typ)
		colValues[i] = val.Interface()
		colEleValues[i] = val.Elem()
	}

	//scan，接收不定长参数，非slice
	if err = rows.Scan(colValues...); err != nil {
		return err
	}

	//再从colEleValues里赋值到r.val
	for i, c := range cs {
		f, _ := r.model.Columns[c]
		r.val.FieldByName(f.FieldName).Set(colEleValues[i])
	}

	return nil
}
