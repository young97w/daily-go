package message

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnDecodeResp(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
	}{
		{
			name: "full",
			resp: &Response{
				Version:    6,
				Compressor: 6,
				Serializer: 61,
				Error:      []byte("error1"),
				Data:       []byte("hello"),
			},
		},
		{
			name: "no err",
			resp: &Response{
				Version:    6,
				Compressor: 6,
				Serializer: 61,
				Data:       []byte("hello"),
			},
		},
		{
			name: "no data",
			resp: &Response{
				Version:    6,
				Compressor: 6,
				Serializer: 61,
				Error:      []byte("error1"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ecd := EncodeResp(tc.resp)
			resp := DecodeResp(ecd)
			assert.Equal(t, tc.resp, resp)
			fmt.Println(string(resp.Data))
		})
	}
}
