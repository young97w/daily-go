package message

import (
	"bytes"
	"encoding/binary"
)

type Request struct {
	HeadLength uint32
	BodyLength uint32
	RequestId  uint32
	Version    uint8
	Compressor uint8
	Serializer uint8

	ServiceName string
	// \n
	MethodName string
	// \n
	//metadata和data都可以省略
	MetaData map[string]string // '/r'分割key value，'/n'分割整个kv
	Data     []byte
}

func EncodeReq(req *Request) []byte {
	req.CalcLength()
	data := make([]byte, req.HeadLength+req.BodyLength+1) //分割用一个/n
	//开始写入头部
	binary.BigEndian.PutUint32(data[:4], req.HeadLength)
	binary.BigEndian.PutUint32(data[4:8], req.BodyLength)
	binary.BigEndian.PutUint32(data[8:12], req.RequestId)
	data[12] = req.Version
	data[13] = req.Compressor
	data[14] = req.Serializer
	//换成cur浅拷贝
	cur := data[15:]
	copy(cur, req.ServiceName)
	cur[len(req.ServiceName)] = '\n'
	cur = cur[len(req.ServiceName)+1:]
	copy(cur, req.MethodName)
	cur = cur[len(req.MethodName):]
	//封装metadata
	cur[0] = '\n'
	cur = cur[1:]
	for k, v := range req.MetaData {
		copy(cur[:len(k)], k)
		cur[len(k)] = '\r'
		copy(cur[len(k)+1:], v)
		cur[len(k)+len(v)+1] = '\n'
		cur = cur[len(k)+len(v)+2:]
	}
	copy(cur, req.Data)
	return data
}

func DecodeReq(data []byte) *Request {
	req := &Request{}
	req.HeadLength = binary.BigEndian.Uint32(data[:4])
	req.BodyLength = binary.BigEndian.Uint32(data[4:8])
	req.RequestId = binary.BigEndian.Uint32(data[8:12])
	req.Version = data[12]
	req.Compressor = data[13]
	req.Serializer = data[14]

	header := data[15:]
	idx := bytes.IndexByte(header, '\n')
	req.ServiceName = string(header[:idx])
	header = header[idx+1:]
	idx = bytes.IndexByte(header, '\n')
	req.MethodName = string(header[:idx])
	header = header[idx+1:]

	//开始解析metadata
	idx = bytes.IndexByte(header, '\n')
	if idx != -1 {
		//如果有metadata
		m := make(map[string]string, 2)
		for idx != -1 {
			pairs := header[:idx]
			sep := bytes.IndexByte(pairs, '\r')
			k := string(pairs[:sep])
			v := string(pairs[sep+1:])
			m[k] = v
			header = header[idx+1:] //不会因为越界报错
			idx = bytes.IndexByte(header, '\n')
		}
		req.MetaData = m
	}
	req.Data = header
	return req
}

func (req *Request) CalcLength() {
	l := 14 + len(req.ServiceName) + 1 + len(req.MethodName) + 1
	for k, v := range req.MetaData {
		l = l + len(k) + 1 + len(v) + 1
	}
	req.HeadLength = uint32(l)
	req.BodyLength = uint32(len(req.Data))
}
