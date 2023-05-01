package v1

import (
	"context"
	"fmt"
	"geektime/ORM/internal/errs"
	"geektime/ORM/v12/internal/model"
	"reflect"
	"strings"
)

//Selector 构造select 语句
type Selector[T any] struct {
	sb      strings.Builder
	args    []any
	table   string
	columns []Selectable
	where   []Predicate
	groupBy []Selectable
	having  []Predicate
	orderBy []Predicate
	offset  int
	limit   int
	model   *model.Model

	db *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		args:  make([]any, 0, 4),
		model: &model.Model{},
		db:    db,
	}
}

// Selectable 指定可以作为选择的表达式
type Selectable interface {
	selectable()
}

//From 指定表名
func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

func (s *Selector[T]) GroupBy(cols ...Selectable) *Selector[T] {
	s.groupBy = cols
	return s
}

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	s.having = ps
	return s
}

func (s *Selector[T]) OrderBy(ps ...Predicate) *Selector[T] {
	s.orderBy = ps
	return s
}

func (s *Selector[T]) Limit(limit, offset int) *Selector[T] {
	s.offset = offset
	s.limit = limit
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//接收db 使用db获取数据 处理结果集
	//先build
	s.limit = 1
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
	model, err := s.db.R.Get(t)
	if err != nil {
		return nil, err
	}
	//新建creator
	creator := s.db.valCreator(t, model) //valuer.NewUnsafeValue(t, model)
	creator.SetColumns(rows)
	return t, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("")
}

func (s *Selector[T]) Build() (*Query, error) {
	var err error
	s.model, err = s.db.R.Get(new(T))
	if err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT ")
	//buildColumn
	err = s.buildColumns()
	if err != nil {
		return nil, err
	}

	s.sb.WriteString(" FROM ")

	if s.table == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.table)
	}

	//构建where
	//注意清空 as
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}
		err = s.buildExpression(p)
		if err != nil {
			return nil, err
		}
	}

	// 构建 GROUP BY
	if len(s.groupBy) > 0 {
		err = s.buildGroupBy()
		if err != nil {
			return nil, err
		}
	}

	//构建 HAVING
	if len(s.having) > 0 {
		err = s.buildHaving()
		if err != nil {
			return nil, err
		}
	}

	//构建 ORDER BY
	if len(s.orderBy) > 0 {
		err = s.buildOrderBy()
		if err != nil {
			return nil, err
		}
	}

	//limit offset
	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.addArgs(s.limit)
	}

	if s.offset > 0 && s.limit > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.addArgs(s.offset)
	}

	//结尾
	s.sb.WriteByte(';')
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.sb.WriteByte('*')
		return nil
	}

	return s.buildGenericCols(s.columns)
}

func (s *Selector[T]) buildGroupBy() error {
	if len(s.groupBy) == 0 {
		return nil
	}

	s.sb.WriteByte(' ')
	s.sb.WriteString("GROUP BY")
	s.sb.WriteByte(' ')
	return s.buildGenericCols(s.groupBy)
}

func (s *Selector[T]) buildHaving() error {
	if len(s.having) == 0 {
		return nil
	}

	s.sb.WriteByte(' ')
	s.sb.WriteString("HAVING")
	s.sb.WriteByte(' ')
	e := s.having[0]

	for i := 1; i < len(s.having); i++ {
		e = e.And(s.having[i])
	}
	return s.buildExpression(e)
}

func (s *Selector[T]) buildOrderBy() error {
	if len(s.orderBy) == 0 {
		return nil
	}

	s.sb.WriteByte(' ')
	s.sb.WriteString("ORDER BY")
	s.sb.WriteByte(' ')

	for i, o := range s.orderBy {
		if i > 0 {
			s.sb.WriteByte(',')
		}
		err := s.buildColumn(o.left.(Column))
		if err != nil {
			return err
		}
		s.sb.WriteByte(' ')
		s.sb.WriteString(o.op.String())
	}
	return nil
}

func (s *Selector[T]) buildGenericCols(cols []Selectable) error {
	for i, column := range cols {
		if i > 0 {
			s.sb.WriteByte(',')
		}

		switch exp := column.(type) {
		case Column:
			err := s.buildColumn(exp)
			if err != nil {
				return err
			}
		case Aggregate:
			s.sb.WriteString(exp.fn)
			s.sb.WriteByte('(')
			err := s.buildColumn(Column{name: exp.col})
			if err != nil {
				return err
			}
			s.sb.WriteByte(')')
			// add alias
			if exp.alias != "" {
				s.sb.WriteString(" AS `")
				s.sb.WriteString(exp.alias)
				s.sb.WriteString("`")
			}
		case RawExpr:
			s.sb.WriteString(exp.raw)
			s.addArgs(exp.args...)
			// add alias
			if exp.alias != "" {
				s.sb.WriteString(" AS `")
				s.sb.WriteString(exp.alias)
				s.sb.WriteString("`")
			}
		}
	}
	return nil
}

func (s *Selector[T]) addArgs(vals ...any) error {
	if len(vals) == 0 {
		return nil
	}

	s.args = append(s.args, vals...)
	return nil
}

func (s *Selector[T]) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}

	switch exp := e.(type) {
	case Column:
		exp.alias = ""
		return s.buildColumn(exp)
	case value:
		s.sb.WriteByte('?')
		s.addArgs(exp.val)
	case Predicate:
		exp.alias = ""
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
		if exp.op.String() != "" {
			s.sb.WriteByte(' ')
			s.sb.WriteString(exp.op.String())
			s.sb.WriteByte(' ')
		}
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
	case RawExpr:
		s.sb.WriteByte('(')
		s.sb.WriteString(exp.raw)
		s.addArgs(exp.args...)
		s.sb.WriteByte(')')
	case Aggregate:
		s.sb.WriteString(exp.fn)
		s.sb.WriteByte('(')
		err := s.buildColumn(Column{name: exp.col})
		if err != nil {
			return err
		}
		s.sb.WriteByte(')')
	default:
		return fmt.Errorf("orm: 不支持的表达式 %v", exp)
	}
	return nil
}

func (s *Selector[T]) Get1(ctx context.Context) (*T, error) {
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

	vals := make([]any, len(cols))
	for i, col := range cols {
		f, ok := s.model.Columns[col]
		if !ok {
			return nil, errs.NewErrUnknownColumn(col)
		}
		val := reflect.New(f.Typ)
		vals[i] = val.Interface()
	}

	if err := rows.Scan(vals...); err != nil {
		return nil, err
	}

	//把vals的值给t
	val := reflect.ValueOf(t).Elem()
	for i, col := range cols {
		f, _ := s.model.Columns[col]
		val.FieldByName(f.FieldName).Set(reflect.ValueOf(vals[i]).Elem())
	}
	return t, nil
}
