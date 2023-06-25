package rpc

import (
	"context"
	"errors"
	__gen "geektime/micro/proto/.gen"
	"geektime/micro/v2/serialize/json"
	"geektime/micro/v2/serialize/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestE2E(t *testing.T) {
	// server
	s := NewServer()
	s.RegisterSerializers(&json.Serializer{})
	s.RegisterService(&UserServiceServer{})
	go func() {
		s.Start("tcp", ":8081")
	}()
	time.Sleep(time.Second * 2)
	// client
	userClient := &UserService{}
	c, err := NewClient(":8081")
	require.NoError(t, err)
	err = c.InitService(userClient)
	require.NoError(t, err)
	res, err := userClient.GetById(context.Background(), &GetByIdReq{Id: 123})
	require.NoError(t, err)
	assert.Equal(t, "123", res.Msg)
}

func TestE2EProto(t *testing.T) {
	// server
	s := NewServer()
	s.RegisterService(&UserServiceServer{})
	go func() {
		s.Start("tcp", ":8081")
	}()
	time.Sleep(time.Second * 2)
	// client
	userClient := &UserService{}

	c, err := NewClient(":8081", ClientWithSerializerOpt(&proto.Serializer{}))
	require.NoError(t, err)
	err = c.InitService(userClient)
	require.NoError(t, err)
	res, err := userClient.GetByIdProto(context.Background(), &__gen.GetByIdReq{Id: 666})
	require.NoError(t, err)
	assert.Equal(t, "666", res.User.Name)
}

func TestCtxWithOneway(t *testing.T) {
	s := NewServer()
	s.RegisterService(&UserService{})
	go func() {
		s.Start("tcp", ":8081")
	}()
	time.Sleep(time.Second * 2)

	userClient := &UserService{}
	c, err := NewClient(":8081")
	require.NoError(t, err)
	err = c.InitService(userClient)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		req     *GetByIdReq
		ctx     context.Context
		wantRes *GetByIdResp
		wantErr error
	}{
		{
			name: "oneway",
			req: &GetByIdReq{
				Id: 123,
			},
			ctx: func() context.Context {
				return CtxWithOneway(context.Background())
			}(),
			wantRes: &GetByIdResp{},
			wantErr: errors.New("micro: 这是一个 oneway 调用，你不应该处理任何结果"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, er := userClient.GetById(tc.ctx, tc.req)
			assert.Equal(t, tc.wantErr, er)
			if er != nil {
				return
			}
			assert.Equal(t, tc.wantRes, resp)
		})
	}

}
