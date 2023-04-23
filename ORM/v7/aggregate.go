package v1

type Aggregate struct {
	fn    string
	col   string
	alias string
}

var _ Selectable = Aggregate{}

func (a Aggregate) selectable() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		col:   a.col,
		alias: alias,
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
