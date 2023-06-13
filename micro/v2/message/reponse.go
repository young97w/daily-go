package message

import (
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
	binary.BigEndian.PutUint32(bs[4:8], resp.BodyLength)
	binary.BigEndian.PutUint32(bs[8:12], resp.RequestID)

	bs[12] = resp.Version
	bs[13] = resp.Compressor
	bs[14] = resp.Serializer
	// write error
	if resp.HeadLength > 15 {
		copy(bs[15:resp.HeadLength], resp.Error)
	}
	if resp.BodyLength > 0 {
		copy(bs[resp.HeadLength:], resp.Data)
	}
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
	if resp.HeadLength > 15 {
		resp.Error = data[15:resp.HeadLength]
	}
	if resp.BodyLength > 0 {
		resp.Data = data[resp.HeadLength:]
	}
	return resp
}

func (resp *Response) CalcLength() {
	resp.HeadLength = uint32(15 + len(resp.Error))
	resp.BodyLength = uint32(len(resp.Data))
}
