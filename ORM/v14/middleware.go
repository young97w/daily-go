package v1

import (
	"context"
	"geektime/ORM/v14/model"
)

type QueryContext struct {
	Type    string
	Builder QueryBuilder
	Model   *model.Model
}

type QueryResult struct {
	Result any
	Err    error
}

type HandleFunc func(ctx context.Context, qc *QueryContext) *QueryResult
type Middleware func(next HandleFunc) HandleFunc
