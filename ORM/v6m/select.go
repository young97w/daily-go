package v1

import (
	"context"
	"fmt"
	"geektime/ORM/internal/errs"
	"reflect"
	"strings"
)

//Selector 构造select 语句
type Selector[T any] struct {
	sb    strings.Builder
	args  []any
	table string
	where []Predicate
	model *Model

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		args:  make([]any, 0, 4),
		model: &Model{},
		db:    db,
	}
}

//From 指定表名
func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//接收db 使用db获取数据 处理结果集
	//先build
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	//处理结果集
	for !rows.Next() {
		return nil, errs.ErrNoRows
	}

	t := new(T)
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	colValues := make([]any, len(cols))
	colEleValues := make([]reflect.Value, len(cols))
	for i, col := range cols {
		f, ok := s.model.columns[col]
		if !ok {
			return nil, errs.NewErrUnknownColumn(col)
		}
		val := reflect.New(f.typ)
		colValues[i] = val.Interface()
		colEleValues[i] = val.Elem()
	}

	if err := rows.Scan(colValues...); err != nil {
		return nil, err
	}

	//把vals的值给t
	val := reflect.ValueOf(t).Elem()
	for i, col := range cols {
		f, _ := s.model.columns[col]
		fd := val.FieldByName(f.fieldName)
		fd.Set(colEleValues[i])
	}
	return t, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("")
}

func (s *Selector[T]) Build() (*Query, error) {
	var err error
	s.model, err = s.db.r.Get(new(T))
	if err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT * FROM ")

	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.tableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}

	//构建where
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		err := s.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}
	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}

	switch exp := e.(type) {
	case Column:
		s.sb.WriteByte('`')
		//s.sb.WriteString(exp.name)
		//校验column
		fd, ok := s.model.fields[exp.name]
		if !ok {
			return errs.NewErrUnknownColumn(exp.name)
		}
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
	case value:
		s.sb.WriteByte('?')
		s.args = append(s.args, exp.val)
	case Predicate:
		//build left
		_, lp := exp.left.(Predicate)
		if lp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			s.sb.WriteByte(')')
		}
		//build operator
		s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')
		//build right
		_, rp := exp.right.(Predicate)
		if rp {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			s.sb.WriteByte(')')
		}
	default:
		return fmt.Errorf("orm: 不支持的表达式 %v", exp)
	}
	return nil
}
