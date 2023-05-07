package querylog

import (
	"context"
	v1 "geektime/ORM/v14"
)

type MiddleBuilder struct {
	logFunc func(sql string, args []any)
}

func NewMiddlewareBuilder() *MiddleBuilder {
	return &MiddleBuilder{logFunc: func(sql string, args []any) {
	}}
}

func (m *MiddleBuilder) LogFunc(logFunc func(sql string, args []any)) *MiddleBuilder {
	m.logFunc = logFunc
	return m
}

func (m *MiddleBuilder) Build() v1.Middleware {
	return func(next v1.HandleFunc) v1.HandleFunc {
		return func(ctx context.Context, qc *v1.QueryContext) *v1.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &v1.QueryResult{Err: err}
			}

			m.logFunc(q.SQL, q.Args)
			return next(ctx, qc)
		}
	}
}
