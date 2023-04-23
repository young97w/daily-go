package v1

import "geektime/ORM/internal/errs"

type Column struct {
	name  string
	alias string
}

var _ Selectable = Column{}

func (c Column) expr()       {}
func (c Column) selectable() {}

type value struct {
	val any
}

func (s *Selector[T]) buildColumn(c Column) error {
	s.sb.WriteByte('`')
	//s.sb.WriteString(exp.name)
	//校验column
	fd, ok := s.model.Fields[c.name]
	if !ok {
		return errs.NewErrUnknownField(c.name)
	}
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')

	// add alias
	if c.alias != "" {
		s.sb.WriteString(" AS `")
		s.sb.WriteString(c.alias)
		s.sb.WriteString("`")
	}
	return nil
}

func (v value) expr() {}

func valueOf(val any) value {
	return value{val: val}
}

func C(field string) Column {
	return Column{name: field}
}

func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

func (c Column) EQ(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: exprOf(val),
	}
}

func (c Column) LT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: exprOf(val),
	}
}

func (c Column) GT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: exprOf(val),
	}
}
