package _reflect

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	testCases := []struct {
		name string

		input any

		wantRes map[string]FuncInfo
		wantErr error
	}{
		{
			// 普通结构体
			name:  "normal struct",
			input: User{},
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name:         "GetAge",
					InputTypes:   []reflect.Type{reflect.TypeOf(User{})},
					OutputTypes:  []reflect.Type{reflect.TypeOf(0)},
					OutputValues: []any{0},
				},
			},
		},
		{
			// 指针
			name:  "pointer",
			input: &User{},
			wantRes: map[string]FuncInfo{
				"GetAge": {
					Name:         "GetAge",
					InputTypes:   []reflect.Type{reflect.TypeOf(&User{})},
					OutputTypes:  []reflect.Type{reflect.TypeOf(0)},
					OutputValues: []any{0},
				},
				"ChangeName": {
					Name:         "ChangeName",
					InputTypes:   []reflect.Type{reflect.TypeOf(&User{}), reflect.TypeOf("")},
					OutputTypes:  []reflect.Type{},
					OutputValues: []any{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFunc(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
