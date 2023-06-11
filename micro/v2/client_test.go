package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func retRespData() []byte {
	resp := &GetByIdResp{Msg: "hello"}
	res, _ := json.Marshal(resp)
	return res
}

func Test_setFuncField(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) Proxy

		service Service
		wantErr error
	}{
		{
			name:    "nil",
			service: nil,
			mock: func(ctrl *gomock.Controller) Proxy {
				return NewMockProxy(ctrl)
			},
			wantErr: errors.New("rpc: 不支持nil"),
		},
		{
			name:    "no pointer",
			service: UserService{},
			mock: func(ctrl *gomock.Controller) Proxy {
				return NewMockProxy(ctrl)
			},
			wantErr: errors.New("rpc: 只支持结构体的一级指针"),
		},
		{
			name: "user service",
			mock: func(ctrl *gomock.Controller) Proxy {
				p := NewMockProxy(ctrl)
				p.EXPECT().Invoke(gomock.Any(), gomock.Any()).Return(&Response{Data: retRespData()}, nil)
				return p
			},
			service: &UserService{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			err := setFuncField(tc.service, tc.mock(ctrl))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			resp, err := tc.service.(*UserService).GetById(context.Background(), &GetByIdReq{Id: 123})
			assert.Equal(t, tc.wantErr, err)
			t.Log(resp)
		})
	}
}

func TestInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	p := NewMockProxy(ctrl)
	p.EXPECT().Invoke(gomock.Any(), gomock.Any()).Return(&Response{Data: retRespData()}, nil)
	res, err := p.Invoke(context.Background(), &Request{})
	require.NoError(t, err)
	msg := &GetByIdResp{}
	json.Unmarshal(res.Data, msg)
	assert.Equal(t, &GetByIdResp{Msg: "hello"}, msg)

}
