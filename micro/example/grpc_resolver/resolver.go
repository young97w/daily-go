package grpc_resolver

import "google.golang.org/grpc/resolver"

type Builder struct {
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		cc: cc,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return nil, nil
}

func (b *Builder) Scheme() string {
	//panic("")
	return "registry"
}

type Resolver struct {
	cc resolver.ClientConn
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	// 固定写死 IP 和端口
	// "localhost:8081"
	err := r.cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{
			{
				Addr: "localhost:8081",
			},
		},
	})
	if err != nil {
		r.cc.ReportError(err)
	}
}

func (r *Resolver) Close() {
	//TODO implement me
	panic("implement me")
}
