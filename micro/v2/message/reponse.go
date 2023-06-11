package message

import (
	"bytes"
	"encoding/binary"
)

type Response struct {
	HeadLength uint32
	BodyLength uint32
	RequestID  uint32
	Version    uint8
	Compressor uint8
	Serializer uint8

	Error []byte
	Data  []byte
}

func EncodeResp(resp *Response) []byte {
	resp.CalcLength()
	bs := make([]byte, resp.HeadLength+resp.BodyLength)

	binary.BigEndian.PutUint32(bs[:4], resp.HeadLength)
	cur := bs[4:]
	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]
	binary.BigEndian.PutUint32(cur[:4], resp.RequestID)
	cur = cur[4:]
	cur[0] = resp.Version
	cur[1] = resp.Compressor
	cur[2] = resp.Serializer
	cur = cur[3:]
	// write error
	if len(resp.Error) > 0 {
		copy(cur, resp.Error)
		cur[len(resp.Error)] = '\n'
	}
	cur = cur[len(resp.Error)+1:]
	copy(cur, resp.Data)
	return bs
}

func DecodeResp(data []byte) *Response {
	resp := &Response{}
	resp.HeadLength = binary.BigEndian.Uint32(data[:4])
	resp.BodyLength = binary.BigEndian.Uint32(data[4:8])
	resp.RequestID = binary.BigEndian.Uint32(data[8:12])
	resp.Version = data[12]
	resp.Compressor = data[13]
	resp.Serializer = data[14]
	cur := data[16:]
	idx := bytes.IndexByte(cur, '\n')
	if idx != -1 {
		resp.Error = cur[:idx]
		cur = cur[idx+1:]
	}
	resp.Data = cur
	return resp
}

func (resp *Response) CalcLength() {
	resp.HeadLength = 15
	resp.BodyLength = uint32(len(resp.Data) + len(resp.Error) + 1)
}
