package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"geektime/micro/v2/message"
	"net"
	"reflect"
)

type Server struct {
	services map[string]ReflectionStub
}

type ReflectionStub struct {
	s     Service
	value reflect.Value
}

func NewServer() *Server {
	return &Server{services: make(map[string]ReflectionStub, 16)}
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = ReflectionStub{
		s:     service,
		value: reflect.ValueOf(service),
	}
}

func (s *Server) Start(network string, addr string) error {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if er := s.handleConn(conn); er != nil {
				_ = conn.Close()
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	for {
		reqBS, err := ReadMsg(conn)
		if err != nil {
			return err
		}
		req := message.DecodeReq(reqBS)
		if err != nil {
			return err
		}
		resp, err := s.Invoke(context.Background(), req)
		if err != nil {
			// handle error
			resp.Error = []byte(err.Error())
		}
		resData := message.EncodeResp(resp)
		_, err = conn.Write(resData)
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	//根据service name 拿到对应的结构体
	refStub, ok := s.services[req.ServiceName]
	resp := &message.Response{
		RequestID:  req.RequestId,
		Version:    req.Version,
		Compressor: req.Compressor,
		Serializer: req.Serializer,
	}
	if !ok {
		return nil, errors.New("调用的服务不存在")
	}
	respData, err := refStub.invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	resp.Data = respData
	return resp, nil
}

func (s *ReflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	method := s.value.MethodByName(req.MethodName)
	in := make([]reflect.Value, 2)
	//第一个参数
	in[0] = reflect.ValueOf(context.Background())
	inReq := reflect.New(method.Type().In(1).Elem())
	err := json.Unmarshal(req.Data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	res := method.Call(in)
	if res[1].Interface() != nil {
		return nil, res[1].Interface().(error)
	}
	return json.Marshal(res[0].Interface())
}
