package micro

import (
	"context"
	"geektime/micro/registry"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ServerOption func(s *Server)

type Server struct {
	name            string
	registry        registry.Registry
	registryTimeout time.Duration
	*grpc.Server
	listener net.Listener
	weight   uint32
	group    string
}

func ServerWithWeight(weight uint32) ServerOption {
	return func(s *Server) {
		s.weight = weight
	}
}

func ServerWithRegister(r registry.Registry) ServerOption {
	return func(s *Server) {
		s.registry = r
	}
}

func NewServer(name string, opts ...ServerOption) (*Server, error) {
	s := &Server{
		name:            name,
		registryTimeout: time.Second * 3,
		Server:          grpc.NewServer(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener

	if s.registry != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.registryTimeout)
		defer cancel()
		er := s.registry.Register(ctx, registry.ServiceInstance{
			Name:    s.name,
			Address: addr,
			Weight:  s.weight,
			Group:   s.group,
		})
		if er != nil {
			return er
		}

	}
	err = s.Serve(s.listener)
	return err
}

func (s *Server) Close() error {
	if s.registry != nil {
		err := s.registry.Close()
		if err != nil {
			return err
		}
	}
	s.GracefulStop()
	return nil
}
