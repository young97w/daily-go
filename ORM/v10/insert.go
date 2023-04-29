package v1

import (
	"reflect"
	"strings"
)

type Inserter[T any] struct {
	sb     strings.Builder
	db     *DB
	values []any
	args   []any
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
	}
}

func (i *Inserter[T]) Values(vals ...any) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	i.sb.WriteString("INSERT INTO ")
	m, err := i.db.R.Get(i.values[0])
	if err != nil {
		return nil, err
	}

	i.sb.WriteByte('`')
	i.sb.WriteString(m.TableName)
	i.sb.WriteByte('`')

	//build columns
	i.sb.WriteByte('(')
	for c, col := range m.ColSlice {
		if c > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('`')
		i.sb.WriteString(col.ColName)
		i.sb.WriteByte('`')
	}
	i.sb.WriteByte(')')

	i.sb.WriteString(" VALUES ")

	//先预估下容量
	i.args = make([]any, 0, len(m.ColSlice))
	//开始插入值
	for _, val := range i.values {
		refVal := reflect.ValueOf(val).Elem()
		i.sb.WriteByte('(')
		for c, col := range m.ColSlice {
			if c > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			i.args = append(i.args, refVal.FieldByName(col.FieldName).Interface())
		}
		i.sb.WriteByte(')')
	}

	i.sb.WriteByte(';')
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}
