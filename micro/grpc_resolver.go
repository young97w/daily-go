package micro

import (
	"context"
	"geektime/micro/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"time"
)

type grpcResolverBuilder struct {
	r       registry.Registry
	timeout time.Duration
}

func NewRegistryBuilder(r registry.Registry, timeout time.Duration) (*grpcResolverBuilder, error) {
	return &grpcResolverBuilder{r: r, timeout: timeout}, nil
}

func (b *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &grpcResolver{
		cc:      cc,
		r:       b.r,
		target:  target,
		timeout: b.timeout,
	}
	r.resolve()
	go r.watch()
	return r, nil
}

func (b *grpcResolverBuilder) Scheme() string {
	return "registry"
}

type grpcResolver struct {
	target  resolver.Target
	r       registry.Registry
	cc      resolver.ClientConn
	timeout time.Duration
	close   chan struct{}
}

func (g *grpcResolver) ResolveNow(options resolver.ResolveNowOptions) {
	g.resolve()
}

func (g *grpcResolver) watch() {
	events, err := g.r.Subscribe(g.target.Endpoint)
	if err != nil {
		g.cc.ReportError(err)
	}
	for {
		select {
		case <-events:
			g.resolve()
		case <-g.close:
			return
		}
	}
}

func (g *grpcResolver) resolve() {
	//获取所有实例再更新
	//fmt.Println("resolve ----")
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	instances, err := g.r.ListService(ctx, g.target.Endpoint)
	if err != nil {
		g.cc.ReportError(err)
	}
	address := make([]resolver.Address, len(instances))
	for _, is := range instances {
		address = append(address, resolver.Address{
			Addr:       is.Address,
			Attributes: attributes.New("weight", is.Weight).WithValue("group", is.Group),
		})
	}

	err = g.cc.UpdateState(resolver.State{
		Addresses: address,
	})

	if err != nil {
		g.cc.ReportError(err)
		return
	}
}

func (g *grpcResolver) Close() {
	close(g.close)
}
