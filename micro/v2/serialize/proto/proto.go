package proto

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

type Serializer struct {
}

func (s *Serializer) Code() uint8 {
	return 1
}

func (s *Serializer) Encode(val any) ([]byte, error) {
	m, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("micro: 必须是 proto.Message")
	}
	return proto.Marshal(m)
}

func (s *Serializer) Decode(data []byte, val any) error {
	m, ok := val.(proto.Message)
	if !ok {
		return errors.New("micro: 必须是 proto.Message")
	}
	return proto.Unmarshal(data, m)
}
