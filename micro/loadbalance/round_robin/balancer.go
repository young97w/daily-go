package round_robin

import (
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync/atomic"
)

type Balancer struct {
	index       int32
	connections []balancer.SubConn
	length      int32
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	idx := atomic.AddInt32(&b.index, 1)
	fmt.Println("Picker Index is:", idx, " length is:", b.length)
	c := b.connections[idx%b.length]
	return balancer.PickResult{
		SubConn: c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for conn := range info.ReadySCs {
		connections = append(connections, conn)
	}
	return &Balancer{
		index:       -1,
		connections: connections,
		length:      int32(len(info.ReadySCs)),
	}
}
