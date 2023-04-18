package sql_demo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMAndUm(t *testing.T) {

	mum := MumTest{
		Name: "aaa",
		Age:  0,
	}

	val, _ := json.Marshal(mum)
	fmt.Println(string(val))

	var m1 MumTest
	json.Unmarshal(val, &m1)

	fmt.Println(m1)
}

type MumTest struct {
	Name string
	Age  int
}

func TestJsonColumn_Value(t *testing.T) {
	//new a JsonColumn with type User
	js := JsonColumn[User]{
		Val:   User{Name: "young"},
		Valid: true,
	}

	value, err := js.Value()
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"Name":"young"}`), value)

	js = JsonColumn[User]{}
	value, err = js.Value()
	require.NoError(t, err)
	require.Nil(t, value)

}

func TestJsonColumn_Scan(t *testing.T) {
	testCase := []struct {
		name string
		src  any

		wantVal   User
		wantValid bool
		wantErr   error
	}{
		{
			name: "nil",
		},
		{
			name: "normal",
			src:  `{"Name":"young"}`,

			wantVal:   User{Name: "young"},
			wantValid: true,
		},
		{
			name: "小写字段名称 src",
			src:  `{"name":"young"}`, //still can unmarshal

			wantVal:   User{Name: "young"},
			wantValid: true,
		},
		{
			name: "错误字段名称 src",
			src:  `{"tname":"young"}`, // can't unmarshal

			wantVal:   User{Name: ""},
			wantValid: true,
		},
		{
			name: "错误类型 src",
			src:  18, // can't unmarshal

			wantValid: false,
			wantErr:   errors.New("orm: 不支持类型"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			js := JsonColumn[User]{}

			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantValid, js.Valid)
			assert.Equal(t, tc.wantVal, js.Val)
		})
	}
}

func TestScanTypes(t *testing.T) {
	//slice , string to val
	jsSlice := JsonColumn[[]string]{}
	err := jsSlice.Scan(`["a","b","c"]`)
	assert.NoError(t, err)
	assert.Equal(t, true, jsSlice.Valid)
	assert.Equal(t, []string{"a", "b", "c"}, jsSlice.Val)

	//slice, val to string
	value, err := jsSlice.Value()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`["a","b","c"]`), value)

	//map string to val
	jsMap := JsonColumn[map[string]string]{}
	err = jsMap.Scan(`{"a":"a"}`)
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "a"}, jsMap.Val)

	//map val to string
	value, err = jsMap.Value()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"a":"a"}`), value)
}

type User struct {
	Name string
}
