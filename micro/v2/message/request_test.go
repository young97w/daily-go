package message

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeReq(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "full",
			req: &Request{
				RequestId:   6,
				Version:     6,
				Compressor:  6,
				Serializer:  61,
				ServiceName: "rpc-demo",
				MethodName:  "foo",
				MetaData: map[string]string{
					"trace":   "tracer",
					"timeout": "1",
				},
				Data: []byte("hello"),
			},
		},
		{
			name: "no data",
			req: &Request{
				RequestId:   6,
				Version:     6,
				Compressor:  6,
				Serializer:  61,
				ServiceName: "rpc-demo",
				MethodName:  "foo",
				MetaData: map[string]string{
					"trace":   "tracer",
					"timeout": "1",
				},
				Data: []byte(""),
			},
		},
		{
			name: "no map",
			req: &Request{
				RequestId:   6,
				Version:     6,
				Compressor:  6,
				Serializer:  61,
				ServiceName: "rpc-demo",
				MethodName:  "foo",
				//MetaData: make(map[string]string,2),
				Data: []byte("hello"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ecd := EncodeReq(tc.req)
			req := DecodeReq(ecd)
			assert.Equal(t, tc.req.Data, req.Data)
			fmt.Println(string(req.Data))
		})
	}
}
