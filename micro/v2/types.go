package rpc

import (
	"context"
	"geektime/micro/v2/message"
)

type Service interface {
	Name() string
}

type Proxy interface {
	Invoke(ctx context.Context, req *message.Request) (*message.Response, error)
}

type Request struct {
	ServiceName string
	MethodName  string
	Arg         []byte
}

type Response struct {
	Data []byte
}
