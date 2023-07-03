package route_round_robin

import (
	"geektime/micro/route"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"sync/atomic"
)

type Balancer struct {
	index       int32
	connections []subConn
	length      int32
	filter      route.Filter
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	candidates := make([]subConn, 0, len(b.connections))
	for _, conn := range b.connections {
		// if filter exists
		if b.filter != nil && !b.filter(info, conn.addr) {
			continue
		}
		candidates = append(candidates, conn)
	}

	if len(candidates) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	idx := atomic.AddInt32(&b.index, 1)
	//fmt.Println("Picker Index is:", idx, " length is:", b.length)
	c := candidates[int(idx)%len(candidates)]
	return balancer.PickResult{
		SubConn: c.c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type Builder struct {
	Filter route.Filter
}

func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	connections := make([]subConn, 0, len(info.ReadySCs))
	for conn, connInfo := range info.ReadySCs {
		connections = append(connections, subConn{
			c:    conn,
			addr: connInfo.Address,
		})
	}
	return &Balancer{
		index:       -1,
		connections: connections,
		length:      int32(len(info.ReadySCs)),
		filter:      b.Filter,
	}
}

type subConn struct {
	c    balancer.SubConn
	addr resolver.Address
}
