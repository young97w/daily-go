package v1

type op string

func (o op) String() string {
	return string(o)
}

const (
	opEQ  op = "="
	opLT  op = "<"
	opGT  op = ">"
	opAND op = "AND"
	opOR  op = "OR"
	opNOT op = "NOT"
)

type Predicate struct {
	left  Expression
	op    op
	right Expression
	alias string
}

func (p Predicate) expr() {}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return valueOf(e)
	}
}

func (p Predicate) As(alias string) Predicate {
	return Predicate{
		left:  p.left,
		op:    p.op,
		right: p.right,
		alias: alias,
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAND,
		right: r,
	}
}

func (p Predicate) Or(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOR,
		right: r,
	}
}
