package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	ID() string
}

type Store interface {
	// 生成一个session
	Generate(ctx context.Context, id string) (Session, error)
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

type Propagator interface {
	// Inject inject 将session注入到响应里面
	Inject(id string, writer http.ResponseWriter) error
	// Extract 提取session id
	Extract(req *http.Request) (string, error)
	//Remove 将session id从http.ResponseWriter中删除
	Remove(writer http.ResponseWriter) error
}
