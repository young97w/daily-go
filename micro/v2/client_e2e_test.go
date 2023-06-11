package rpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestE2E(t *testing.T) {
	// server
	s := NewServer()
	s.RegisterService(&UserServiceServer{})
	go func() {
		s.Start("tcp", ":8081")
	}()
	time.Sleep(time.Second * 2)
	// client
	userClient := &UserService{}
	err := InitClientProxy(":8081", userClient)
	require.NoError(t, err)
	res, err := userClient.GetById(context.Background(), &GetByIdReq{Id: 123})
	require.NoError(t, err)
	assert.Equal(t, "空尼基哇", res.Msg)

}
