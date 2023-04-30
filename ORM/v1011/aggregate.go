package v1

type Aggregate struct {
	fn    string
	col   string
	alias string
}

var _ Selectable = Aggregate{}
var _ Expression = Aggregate{}

func (a Aggregate) selectable() {}

func (a Aggregate) expr() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		col:   a.col,
		alias: alias,
	}
}

func (a Aggregate) EQ(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opEQ,
		right: valueOf(val),
	}
}

func (a Aggregate) LT(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opLT,
		right: valueOf(val),
	}
}

func (a Aggregate) GT(val any) Predicate {
	return Predicate{
		left:  a,
		op:    opLT,
		right: valueOf(val),
	}
}

func Max(col string) Aggregate {
	return Aggregate{
		fn:  "MAX",
		col: col,
	}
}

func Min(col string) Aggregate {
	return Aggregate{
		fn:  "MIN",
		col: col,
	}
}

func Avg(col string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		col: col,
	}
}
