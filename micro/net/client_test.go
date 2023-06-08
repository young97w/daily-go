package net

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Send(t *testing.T) {
	s := NewServer("tcp", ":8081")
	go func() {
		err := s.Start()
		require.NoError(t, err)
	}()

	time.Sleep(time.Second * 3)

	c := NewClient("tcp", ":8081")
	res, err := c.Send([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, "hello", string(res))
}
