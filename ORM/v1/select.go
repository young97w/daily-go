package v1

import (
	"reflect"
	"strings"
)

//Selector 构造select 语句
type Selector[T any] struct {
	table string
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}

//From 指定表名
func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T
		sb.WriteString(reflect.TypeOf(t).Name())
	} else {
		sb.WriteString(s.table)
	}
	sb.WriteByte(';')
	return &Query{SQL: sb.String()}, nil
}
