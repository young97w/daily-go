package v1

type Expression interface {
	expr()
}

// RawExpr 代表的是原生表达式
type RawExpr struct {
	raw   string
	args  []any
	alias string
}

func (r RawExpr) As(alias string) RawExpr {
	return RawExpr{
		raw:   r.raw,
		args:  r.args,
		alias: r.alias,
	}
}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

func (r RawExpr) selectable() {}
func (r RawExpr) expr()       {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}
