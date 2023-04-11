package _reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MUser struct {
	Name string
	age  int
}

func TestIterateFields(t *testing.T) {
	testCase := []struct {
		name string

		entity any

		wantErr error
		wantRes map[string]any
	}{
		{
			name: "age zero",
			entity: User{
				Name: "Young",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Young",
				"age":  0,
			},
		},
		{
			name: "pointer",
			entity: &User{
				Name: "Young",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Young",
				"age":  0,
			},
		},
		{
			name:    "basic type",
			entity:  18,
			wantErr: errors.New("不支持类型"),
		},
		{
			name:    "nil",
			entity:  nil,
			wantErr: errors.New("不支持nil"),
		},
		{
			name:    "user nil",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持零值"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFields(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestSetValue(t *testing.T) {
	testCase := []struct {
		name string

		entity   any
		field    string
		newValue any

		wantErr error
		wantRes any
	}{
		{
			name:     "normal struct",
			entity:   &User{},
			field:    "Name",
			newValue: "kk",
			wantRes:  &User{Name: "kk"},
		},
		{
			name:     "unaddressable",
			entity:   User{},
			field:    "Name",
			newValue: "kk",
			wantErr:  errors.New("不可更改字段"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			err := SetValue(tc.entity, tc.field, tc.newValue)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, tc.entity)
		})
	}
}
