package route_round_robin

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/balancer"
	"testing"
)

func TestWeightBalancer_Pick(t *testing.T) {
	b := &WeightBalancer{
		connections: []*weightConn{
			{
				c: SubConn{
					name: "weight-5",
				},
				weight:          5,
				efficientWeight: 5,
				currentWeight:   5,
			},
			{
				c: SubConn{
					name: "weight-4",
				},
				weight:          4,
				efficientWeight: 4,
				currentWeight:   4,
			},
			{
				c: SubConn{
					name: "weight-3",
				},
				weight:          3,
				efficientWeight: 3,
				currentWeight:   3,
			},
		},
	}
	pickRes, err := b.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	assert.Equal(t, "weight-5", pickRes.SubConn.(SubConn).name)

	pickRes, err = b.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	assert.Equal(t, "weight-4", pickRes.SubConn.(SubConn).name)

	pickRes, err = b.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	assert.Equal(t, "weight-3", pickRes.SubConn.(SubConn).name)

	pickRes, err = b.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	assert.Equal(t, "weight-5", pickRes.SubConn.(SubConn).name)

	pickRes, err = b.Pick(balancer.PickInfo{})
	require.NoError(t, err)
	assert.Equal(t, "weight-4", pickRes.SubConn.(SubConn).name)

	pickRes.Done(balancer.DoneInfo{})
	// 断言这里面 efficient weight 是变化了的
}
