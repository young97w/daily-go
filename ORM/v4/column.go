package v1

type Column struct {
	name string
}

func (c Column) expr() {}

type value struct {
	val any
}

func (v value) expr() {}

func valueOf(val any) value {
	return value{val: val}
}

func C(field string) Column {
	return Column{name: field}
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
