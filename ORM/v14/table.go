package v1

type TableReference interface {
	table()
}

type Table struct {
	entity any
	alias  string
}

func (t Table) table() {}

// TableOf new 一个Table的结构体
func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

func (t Table) C(name string) Column {
	return Column{
		name:  name,
		table: t,
	}
}

func (t Table) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "LEFT JOIN",
		right: right,
	}
}

func (t Table) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "JOIN",
		right: right,
	}
}

func (t Table) RightJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "RIGHT JOIN",
		right: right,
	}
}

type JoinBuilder struct {
	left  TableReference
	typ   string
	right TableReference
}

func (j *JoinBuilder) On(ps ...Predicate) Join {
	return Join{
		left:  j.left,
		typ:   j.typ,
		right: j.right,
		on:    ps,
	}
}

func (j *JoinBuilder) Using(cols ...string) Join {
	return Join{
		left:  j.left,
		typ:   j.typ,
		right: j.right,
		using: cols,
	}
}

type Join struct {
	left  TableReference
	typ   string
	right TableReference
	on    []Predicate
	using []string
}

func (j Join) table() {}

func (j Join) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		typ:   "JOIN",
		right: right,
	}
}

func (j Join) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		typ:   "LEFT JOIN",
		right: right,
	}
}

func (j Join) RightJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		typ:   "RIGHT JOIN",
		right: right,
	}
}
