package homework_delete

import (
	"errors"
	"strings"
)

type Deleter[T any] struct {
	tableName string
	model     *model
	sb        *strings.Builder

	//predicates
	where []Predicate
	args  []any
}

func (d *Deleter[T]) Build() (*Query, error) {
	var err error
	d.model, err = parseModel(new(T))
	if err != nil {
		return nil, err
	}
	d.sb = &strings.Builder{}

	//开始build
	d.sb.WriteString("DELETE ")
	d.sb.WriteString("FROM ")

	if d.tableName != "" {
		d.sb.WriteString(d.tableName)
	} else {
		d.sb.WriteByte('`')
		d.sb.WriteString(d.model.tableName)
		d.sb.WriteByte('`')
	}

	//if where
	if len(d.where) > 0 {
		d.sb.WriteString(" WHERE ")
		p := d.where[0]
		//组装predicate
		for i := 1; i < len(d.where); i++ {
			p = p.And(d.where[i])
		}
		//使用predicate 构建表达式
		d.buildExpression(p)
	}

	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) buildExpression(exp Expression) error {
	//build expression 递归构建
	//递归：先构建左边再构建右边
	//类型：column value predicate
	if exp == nil {
		return nil
	}
	switch expr := exp.(type) {
	case Column:
		d.sb.WriteByte('`')
		d.sb.WriteString(expr.name)
		d.sb.WriteByte('`')
	case value:
		d.sb.WriteByte('?')
		d.args = append(d.args, expr.val)
	case Predicate:
		//build left
		_, lp := expr.left.(Predicate)
		if lp {
			d.sb.WriteByte('(')
		}
		//如果是col 或 value 则不需要括号
		d.buildExpression(expr.left)
		if lp {
			d.sb.WriteByte(')')
		}

		//build operator
		d.sb.WriteByte(' ')
		d.sb.WriteString(expr.op.String())
		d.sb.WriteByte(' ')

		//build right
		_, rp := expr.right.(Predicate)
		if rp {
			d.sb.WriteByte('(')
		}
		d.buildExpression(expr.right)
		if rp {
			d.sb.WriteByte(')')
		}
	default:
		return errors.New("orm: 不支持的表达式")

	}

	return nil
}

// From accepts model definition
func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.tableName = table
	return d
}

// Where accepts predicates
func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	if d.where == nil {
		d.where = make([]Predicate, len(predicates))
	}
	d.where = predicates
	return d
}
