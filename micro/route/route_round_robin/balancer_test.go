package route_round_robin

import (
	"google.golang.org/grpc/balancer"
)

//func TestBalancer_Pick(t *testing.T) {
//	testCases := []struct {
//		name string
//		b    *Balancer
//
//		wantErr           error
//		wantSubConn       SubConn
//		wantBalancerIndex int32
//	}{
//		{
//			name: "start",
//			b: &Balancer{
//				index: -1,
//				connections: []balancer.SubConn{
//					SubConn{name: "127.0.0.1:8080"},
//					SubConn{name: "127.0.0.1:8081"},
//				},
//				length: 2,
//			},
//			wantSubConn:       SubConn{name: "127.0.0.1:8080"},
//			wantBalancerIndex: 0,
//		},
//		{
//			name: "end",
//			b: &Balancer{
//				index: 1,
//				connections: []balancer.SubConn{
//					SubConn{name: "127.0.0.1:8080"},
//					SubConn{name: "127.0.0.1:8081"},
//				},
//				length: 2,
//			},
//
//			wantSubConn:       SubConn{name: "127.0.0.1:8080"},
//			wantBalancerIndex: 2,
//		},
//		{
//			name: "no connections",
//			b: &Balancer{
//				index:       -1,
//				connections: []balancer.SubConn{},
//			},
//			wantErr: balancer.ErrNoSubConnAvailable,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			res, err := tc.b.Pick(balancer.PickInfo{})
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//
//			assert.Equal(t, tc.wantSubConn.name, res.SubConn.(SubConn).name)
//			assert.NotNil(t, res.Done)
//			assert.Equal(t, tc.wantBalancerIndex, tc.b.index)
//		})
//	}
//}

type SubConn struct {
	name string
	balancer.SubConn
}
