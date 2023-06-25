package rpc

import (
	"context"
	"errors"
	"geektime/micro/v2/message"
	"geektime/micro/v2/serialize"
	"geektime/micro/v2/serialize/json"
	"geektime/micro/v2/serialize/proto"
	"net"
	"reflect"
)

type Server struct {
	services    map[string]ReflectionStub
	serializers map[uint8]serialize.Serializer
}

type ReflectionStub struct {
	s           Service
	value       reflect.Value
	serializers map[uint8]serialize.Serializer
}

func NewServer() *Server {
	s := &Server{services: make(map[string]ReflectionStub, 16), serializers: make(map[uint8]serialize.Serializer, 2)}
	s.RegisterSerializers(&json.Serializer{}, &proto.Serializer{})
	return s
}

func (s *Server) RegisterSerializers(sl ...serialize.Serializer) {
	for _, serializer := range sl {
		s.serializers[serializer.Code()] = serializer
	}
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = ReflectionStub{
		s:           service,
		value:       reflect.ValueOf(service),
		serializers: s.serializers,
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
	service, ok := s.services[req.ServiceName]
	resp := &message.Response{
		RequestID:  req.RequestId,
		Version:    req.Version,
		Compressor: req.Compressor,
		Serializer: req.Serializer,
	}
	if !ok {
		return nil, errors.New("调用的服务不存在")
	}
	if isOneway(ctx) {
		go func() {
			_, _ = service.invoke(ctx, req)
		}()
		return nil, errors.New("micro: 微服务服务端 oneway 请求")
	}
	respData, err := service.invoke(ctx, req)
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
	sl, ok := s.serializers[req.Serializer]
	if !ok {
		return nil, errors.New("micro: 不支持的序列化协议")
	}
	err := sl.Decode(req.Data, inReq.Interface())
	if err != nil {
		return nil, err
	}
	in[1] = inReq
	res := method.Call(in)
	if res[1].Interface() != nil {
		return nil, res[1].Interface().(error)
	}
	return sl.Encode(res[0].Interface())
}
