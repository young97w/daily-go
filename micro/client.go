package micro

import (
	"context"
	"fmt"
	"geektime/micro/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"time"
)

type Client struct {
	insecure bool
	r        registry.Registry
	timeout  time.Duration
	balancer balancer.Builder
}

type ClientOption func(client *Client)

func NewClient(opts ...ClientOption) (*Client, error) {
	res := &Client{}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func ClientInsecure() ClientOption {
	return func(client *Client) {
		client.insecure = true
	}
}

func ClientWithRegistry(r registry.Registry, timeout time.Duration) ClientOption {
	return func(client *Client) {
		client.r = r
		client.timeout = timeout
	}
}

func ClientWithPickerBuilder(name string, b base.PickerBuilder) ClientOption {
	return func(c *Client) {
		builder := base.NewBalancerBuilder(name, b, base.Config{HealthCheck: true})
		balancer.Register(builder)
		c.balancer = builder
	}
}

func (c *Client) Dial(ctx context.Context, service string, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	if c.r != nil {
		rb, err := NewRegistryBuilder(c.r, c.timeout)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithResolvers(rb))
	}

	if c.insecure {
		opts = append(opts, grpc.WithInsecure())
	}

	if c.balancer != nil {
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, c.balancer.Name())))
	}

	if len(dialOptions) > 0 {
		opts = append(opts, dialOptions...)
	}

	//dial
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("registry:///%s", service), opts...)
	return conn, err
}
