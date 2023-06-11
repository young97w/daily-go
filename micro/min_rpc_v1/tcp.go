package rpc

import (
	"encoding/binary"
	"net"
)

// 长度字段使用的字节数量
const numOfLengthBytes = 8

func ReadMsg(conn net.Conn) ([]byte, error) {
	lenBS := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBS)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint64(lenBS)
	data := make([]byte, length)
	_, err = conn.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EncodeMsg(data []byte) []byte {
	reqLen := len(data)
	res := make([]byte, reqLen+numOfLengthBytes)
	binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(reqLen))
	copy(res[numOfLengthBytes:], data)
	return res
}
