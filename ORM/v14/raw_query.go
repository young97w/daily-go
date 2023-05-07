package v1

import (
	"context"
	"database/sql"
)

type RawQuerier[T any] struct {
	core
	sess Session
	sql  string
	args []any
}

func RawQuery[T any](sess Session, sql string, args ...any) *RawQuerier[T] {
	c := sess.getCore()
	return &RawQuerier[T]{
		core: c,
		sql:  sql,
		args: args,
		sess: sess,
	}
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

// Exec 对应 UPDATE DELETE INSERT
func (r *RawQuerier[T]) Exec(ctx context.Context) (sql.Result, error) {
	m, err := r.R.Get(new(T))
	if err != nil {
		return nil, err
	}
	var handler HandleFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return execHandler[T](r.sess, ctx, qc)
	}
	result := handler(ctx, &QueryContext{
		Type:    "RAW",
		Builder: r,
		Model:   m,
	})
	return result.Result.(sql.Result), nil
}

// Get 对应SELECT语句
func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	var handler HandleFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](r.sess, r.core, ctx, qc)
	}
	m, err := r.R.Get(new(T))
	if err != nil {
		return nil, err
	}
	result := handler(ctx, &QueryContext{
		Type:    "RAW",
		Builder: r,
		Model:   m,
	})

	return result.Result.(*T), nil
}

func (r *RawQuerier[T]) GetMulti(tx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
