package route

import (
	"context"
	__gen "geektime/micro/proto/.gen"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestServer(t *testing.T) {
	//etcdClient, err := clientv3.New(clientv3.Config{
	//	Endpoints: []string{"localhost:2379"},
	//})
	//require.NoError(t, err)
	//r, err := etcd.NewRegistry(etcdClient)
	//require.NoError(t, err)
	//
	//var eg errgroup.Group
	//for i := 0; i < 3; i++ {
	//	micro.NewServer()
	//}
}

type UserServiceServer struct {
	group string
	__gen.UnimplementedUserServiceServer
}

func (s UserServiceServer) GetById(ctx context.Context, req *__gen.GetByIdReq) (*__gen.GetByIdResp, error) {
	//go func() {
	// 转异步
	//	fmt.Println(s.group)
	//	// 做一些事情
	//}()
	// 返回一个 202
	//fmt.Println(s.group)
	return &__gen.GetByIdResp{
		User: &__gen.User{
			Name: "hello, world",
		},
	}, nil
}

func TestEtcd(t *testing.T) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)

	_, err = c.Put(context.Background(), "key", "val")
	require.NoError(t, err)

	res, err := c.Get(context.Background(), "", clientv3.WithPrefix())
	require.NoError(t, err)
	t.Log(res.Kvs)

	_, err = c.Delete(context.Background(), "key")
	require.NoError(t, err)
}
