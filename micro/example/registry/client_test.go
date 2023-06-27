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

func TestClient(t *testing.T) {
	etcdClient, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	require.NoError(t, err)
	r, err := etcd.NewRegistry(etcdClient)
	require.NoError(t, err)

	client, err := micro.NewClient(micro.ClientWithRegistry(r, time.Second*2), micro.ClientInsecure())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cc, err := client.Dial(ctx, "user-service")
	require.NoError(t, err)

	uc := __gen.NewUserServiceClient(cc)
	resp, err := uc.GetById(ctx, &__gen.GetByIdReq{Id: 123})
	require.NoError(t, err)
	t.Log(resp)
}
