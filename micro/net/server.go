package net

import (
	"encoding/binary"
	"net"
)

const lengthBytes = 8

type Server struct {
	Network string
	Addr    string
}

func NewServer(network, addr string) *Server {
	return &Server{
		Network: network,
		Addr:    addr,
	}
}

// Start 监听网络
func (s *Server) Start() error {
	listener, err := net.Listen(s.Network, s.Addr)
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
	//保持连接
	for {
		lenBS := make([]byte, lengthBytes)
		_, err := conn.Read(lenBS)
		if err != nil {
			return err
		}
		reqLength := binary.BigEndian.Uint64(lenBS)
		reqData := make([]byte, reqLength)
		_, err = conn.Read(reqData)
		if err != nil {
			return err
		}
		respData := make([]byte, lengthBytes+reqLength)
		binary.BigEndian.PutUint64(respData[:lengthBytes], reqLength)
		copy(respData[lengthBytes:], reqData)

		_, err = conn.Write(respData)
		if err != nil {
			return err
		}
	}
}
