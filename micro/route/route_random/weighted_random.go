package route_random

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math/rand"
)

type WeightedBalancer struct {
	connections []*weightConn
	totalWeight uint32
}

type weightConn struct {
	c      balancer.SubConn
	weight uint32
}

func (b *WeightedBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(b.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	tgt := rand.Intn(int(b.totalWeight + 1))
	//fmt.Println("Total weight:", b.totalWeight)
	//fmt.Println("tgt weight:", tgt)
	var res *weightConn
	for _, c := range b.connections {
		tgt -= int(c.weight)
		//fmt.Println("Weight ", c.weight)
		if tgt <= 0 {
			res = c
			break
		}
	}
	//if res == nil {
	//	res = b.connections[0]
	//}
	//fmt.Println("Picker tgt is:", tgt, " length is:", len(b.connections))
	return balancer.PickResult{
		SubConn: res.c,
		Done: func(info balancer.DoneInfo) {

		},
	}, nil
}

type WeightedBalanceBuilder struct {
}

func (b *WeightedBalanceBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	cs := make([]*weightConn, 0, len(info.ReadySCs))
	var totalWeight uint32
	for conn, connInfo := range info.ReadySCs {
		weight := connInfo.Address.Attributes.Value("weight").(uint32)
		totalWeight += weight
		cs = append(cs, &weightConn{c: conn, weight: weight})
	}

	return &WeightedBalancer{
		connections: cs,
		totalWeight: totalWeight,
	}
}
