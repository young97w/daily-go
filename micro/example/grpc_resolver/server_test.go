package grpc_resolver

import (
	"context"
	__gen "geektime/micro/proto/.gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	us := &Server{}
	server := grpc.NewServer()
	__gen.RegisterUserServiceServer(server, us)
	l, err := net.Listen("tcp", ":8081")
	require.NoError(t, err)
	err = server.Serve(l)
	t.Log(err)
}

type Server struct {
	__gen.UnimplementedUserServiceServer
}

func (s Server) GetById(ctx context.Context, req *__gen.GetByIdReq) (*__gen.GetByIdResp, error) {
	return &__gen.GetByIdResp{
		User: &__gen.User{Name: "yong"},
	}, nil
}
