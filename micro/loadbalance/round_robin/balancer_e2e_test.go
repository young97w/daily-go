package round_robin

import (
	"context"
	"fmt"
	__gen "geektime/micro/proto/.gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"net"
	"testing"
	"time"
)

func TestBalancer_e2e_Pick(t *testing.T) {
	go func() {
		us := &Server{}
		server := grpc.NewServer()
		__gen.RegisterUserServiceServer(server, us)
		l, err := net.Listen("tcp", ":8081")
		require.NoError(t, err)
		err = server.Serve(l)
		t.Log(err)
	}()

	time.Sleep(time.Second * 3)
	balancer.Register(base.NewBalancerBuilder("DEMO_ROUND_ROBIN", &Builder{}, base.Config{HealthCheck: true}))
	cc, err := grpc.Dial("localhost:8081", grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "DEMO_ROUND_ROBIN"}`))
	require.NoError(t, err)
	client := __gen.NewUserServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	resp, err := client.GetById(ctx, &__gen.GetByIdReq{Id: 13})
	require.NoError(t, err)
	t.Log(resp)
}

type Server struct {
	__gen.UnimplementedUserServiceServer
}

func (s Server) GetById(ctx context.Context, req *__gen.GetByIdReq) (*__gen.GetByIdResp, error) {
	fmt.Println(req)
	return &__gen.GetByIdResp{
		User: &__gen.User{
			Name: "hello, world",
		},
	}, nil
}
