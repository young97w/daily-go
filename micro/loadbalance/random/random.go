package random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

type Balancer struct {
	connections []balancer.SubConn
	length      int
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if b.length == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	idx := rand.Intn(b.length)
	return balancer.PickResult{
		SubConn: b.connections[idx],
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type BalancerBuilder struct {
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	cs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for conn, _ := range info.ReadySCs {
		cs = append(cs, conn)
	}

	return &Balancer{
		connections: cs,
		length:      len(cs),
	}
}
