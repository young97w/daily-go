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
	length := binary.BigEndian.Uint32(lenBS[:4]) + binary.BigEndian.Uint32(lenBS[4:])
	data := make([]byte, length)
	_, err = conn.Read(data[8:])
	if err != nil {
		return nil, err
	}
	copy(data[:8], lenBS)
	return data, nil
}

//func EncodeMsg(data []byte) []byte {
//	reqLen := len(data)
//	res := make([]byte, reqLen+numOfLengthBytes)
//	binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(reqLen))
//	copy(res[numOfLengthBytes:], data)
//	return res
//}
