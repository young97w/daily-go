package net

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	network string
	addr    string
}

func NewClient(network, addr string) *Client {
	return &Client{
		network: network,
		addr:    addr,
	}
}

func (c *Client) Send(data []byte) ([]byte, error) {
	//先连接
	conn, err := net.DialTimeout(c.network, c.addr, time.Second*2)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		return nil, err
	}
	//send req
	reqLen := len(data)
	reqData := make([]byte, reqLen+lengthBytes)
	binary.BigEndian.PutUint64(reqData[:lengthBytes], uint64(reqLen))
	copy(reqData[lengthBytes:], data)
	_, err = conn.Write(reqData)
	if err != nil {
		return nil, err
	}

	//receive response
	lenBS := make([]byte, lengthBytes)
	_, err = conn.Read(lenBS)
	if err != nil {
		return nil, err
	}
	respLength := binary.BigEndian.Uint64(lenBS)
	respData := make([]byte, respLength)
	_, err = conn.Read(respData)
	if err != nil {
		return nil, err
	}
	return respData, nil
}
