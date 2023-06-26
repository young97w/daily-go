package registry

import (
	"context"
	"io"
)

type ServiceInstance struct {
	Name    string
	Address string

	//custom field
	Weight uint32
	Group  string
}

type Event struct {
	Type string
}

type Registry interface {
	Register(ctx context.Context, si ServiceInstance) error
	Unregister(ctx context.Context, si ServiceInstance) error
	ListService(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	Subscribe(serviceName string) (<-chan Event, error)

	io.Closer
}
