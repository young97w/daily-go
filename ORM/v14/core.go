package v1

import (
	"context"
	"geektime/ORM/v14/internal/valuer"
	"geektime/ORM/v14/model"
)

type core struct {
	dialect    Dialect
	valCreator valuer.Creator
	R          *model.Registry
	mdls       []Middleware
}

func getHandler[T any](sess Session, c core, ctx context.Context, qc *QueryContext) *QueryResult {
	//接收db 使用db获取数据 处理结果集
	//先build
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{Err: err}
	}

	rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{Err: err}
	}

	//处理结果集
	for !rows.Next() {
		return nil
	}

	t := new(T)
	model, err := c.R.Get(t)
	if err != nil {
		return &QueryResult{Err: err}
	}
	//新建creator
	creator := c.valCreator(t, model) //valuer.NewUnsafeValue(t, model)
	err = creator.SetColumns(rows)
	return &QueryResult{Result: t, Err: err}
}

func execHandler[T any](sess Session, ctx context.Context, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{Err: err}
	}

	res, err := sess.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{Result: res, Err: err}
}
