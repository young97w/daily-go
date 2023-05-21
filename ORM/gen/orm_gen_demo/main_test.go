package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGen(t *testing.T) {
	src := "testdata/user.go"
	b := &bytes.Buffer{}
	err := gen(b, src)
	require.NoError(t, err)
	fmt.Println(b)
}
