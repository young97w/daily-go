package rpc

import (
	"context"
	__gen "geektime/micro/proto/.gen"
	"log"
	"strconv"
)

type UserService struct {
	GetById      func(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error)
	GetByIdProto func(ctx context.Context, req *__gen.GetByIdReq) (resp *__gen.GetByIdResp, err error)
}

func (u UserService) Name() string {
	return "user-service"
}

type GetByIdReq struct {
	Id int
}

type GetByIdResp struct {
	Msg string
}

type UserServiceServer struct {
}

func (u *UserServiceServer) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {
	log.Println(req)
	return &GetByIdResp{Msg: strconv.Itoa(req.Id)}, nil
}

func (u *UserServiceServer) GetByIdProto(ctx context.Context, req *__gen.GetByIdReq) (*__gen.GetByIdResp, error) {
	log.Println(req)
	return &__gen.GetByIdResp{User: &__gen.User{
		Name: strconv.Itoa(int(req.Id)),
	}}, nil
}

func (u *UserServiceServer) Name() string {
	return "user-service"
}
