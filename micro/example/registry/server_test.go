package grpc

import (
	"context"
	"geektime/micro"
	__gen "geektime/micro/proto/.gen"
	"geektime/micro/registry/etcd"
	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	registry, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)
	server, err := micro.NewServer("user-service", micro.ServerWithRegister(registry))
	require.NoError(t, err)
	us := &UserServiceServer{}
	__gen.RegisterUserServiceServer(server, us)
	err = server.Start(":8081")
	t.Log(err)
}

type UserServiceServer struct {
	__gen.UnimplementedUserServiceServer
}

func (s UserServiceServer) GetById(ctx context.Context, req *__gen.GetByIdReq) (*__gen.GetByIdResp, error) {
	return &__gen.GetByIdResp{
		User: &__gen.User{Name: "yong"},
	}, nil
}

func TestEtcd(t *testing.T) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	response, err := c.Put(ctx, "key", "val")
	require.NoError(t, err)
	t.Log(response)
	deleteResponse, err := c.Delete(ctx, "key")
	require.NoError(t, err)
	t.Log(deleteResponse)
}
