package route_round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"sync"
)

type weightConn struct {
	mutex           sync.Mutex
	c               balancer.SubConn
	weight          uint32
	currentWeight   uint32
	efficientWeight uint32
}

type WeightBalancer struct {
	connections []*weightConn
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 使用weighted round robin选出 conn
	// 目标conn，current weight 归零
	if len(w.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var res *weightConn
	var totalWeight uint32
	for _, c := range w.connections {
		c.mutex.Lock()
		totalWeight += c.efficientWeight
		c.currentWeight += c.efficientWeight
		if res == nil || res.currentWeight < c.currentWeight {
			res = c
		}
		c.mutex.Unlock()
	}

	res.mutex.Lock()
	res.currentWeight -= totalWeight
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.c,
		Done: func(info balancer.DoneInfo) {
			res.mutex.Lock()
			if info.Err != nil && res.efficientWeight == 0 {
				return
			}
			if info.Err == nil && res.efficientWeight == math.MaxUint32 {
				return
			}
			if info.Err != nil {
				res.efficientWeight--
			} else {
				res.efficientWeight++
			}
			res.mutex.Unlock()
		},
	}, nil
}

type WeightBalanceBuilder struct {
}

func (w *WeightBalanceBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	cs := make([]*weightConn, 0, len(info.ReadySCs))
	for sub, subInfo := range info.ReadySCs {
		weight := subInfo.Address.Attributes.Value("weight").(uint32)
		cs = append(cs, &weightConn{
			c:               sub,
			weight:          weight,
			currentWeight:   weight,
			efficientWeight: weight,
		})
	}

	return &WeightBalancer{connections: cs}
}
