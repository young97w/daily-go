package round_robin

import (
	"context"
	"fmt"
	"geektime/micro"
	"geektime/micro/loadbalance/random"
	__gen "geektime/micro/proto/.gen"
	"geektime/micro/registry/etcd"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"sync"
	"testing"
	"time"
)

func TestBalancer_e2e_Pick(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)

	go func() {
		var eg errgroup.Group
		for i := 0; i < 3; i++ {
			f := func(i int) func() {
				return func() {
					fmt.Println("server:", i)
				}
			}(i)

			us := &Server{f: f}
			server, err2 := micro.NewServer("user-service", micro.ServerWithRegister(r), micro.ServerWithWeight(uint32(4-i)))
			require.NoError(t, err2)
			__gen.RegisterUserServiceServer(server, us)
			require.NoError(t, err2)
			port := fmt.Sprintf(":808%d", i)
			eg.Go(func() error {
				return server.Start(port)
			})
		}
		err := eg.Wait()
		t.Log(err)
	}()

	time.Sleep(time.Second * 3)

	client, err := micro.NewClient(micro.ClientInsecure(), micro.ClientWithRegistry(r, time.Second*2),
		micro.ClientWithPickerBuilder("DEMO_ROUND_ROBIN", &random.WeightedBalanceBuilder{}))
	require.NoError(t, err)

	cc, err := client.Dial(context.Background(), "user-service")
	uc := __gen.NewUserServiceClient(cc)
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		resp, err := uc.GetById(ctx, &__gen.GetByIdReq{Id: int64(i)})
		cancel()
		require.NoError(t, err)
		t.Log(resp)
		//go func() {
		//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		//	resp, err := uc.GetById(ctx, &__gen.GetByIdReq{Id: int64(i)})
		//	defer cancel()
		//	require.NoError(t, err)
		//	t.Log(resp)
		//
		//}()
	}
	wg.Wait()
}

type Server struct {
	__gen.UnimplementedUserServiceServer
	f func()
}

func (s Server) GetById(ctx context.Context, req *__gen.GetByIdReq) (*__gen.GetByIdResp, error) {
	fmt.Println(req)
	s.f()
	return &__gen.GetByIdResp{
		User: &__gen.User{
			Name: "hello, world",
		},
	}, nil
}
