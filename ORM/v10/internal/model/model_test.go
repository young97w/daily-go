package model

import (
	"geektime/ORM/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	testCases := []struct {
		name string

		entity  any
		wantRes *Model
		wantErr error
	}{
		{
			name:   "normal struct",
			entity: testModel{},
			wantRes: &Model{
				TableName: "test_model",
				Fields: map[string]*Field{
					"Name": &Field{ColName: "name"},
					"Age":  &Field{ColName: "age"},
				},
				Columns: map[string]*Field{
					"name": {FieldName: "Name", Typ: reflect.TypeOf("")},
					"age":  {FieldName: "Age", Typ: reflect.TypeOf(18)},
				},
			},
		},
		{
			name:   "with 1 pointer",
			entity: &testModel{},
			wantRes: &Model{
				TableName: "test_model",
				Fields: map[string]*Field{
					"Name": &Field{ColName: "name"},
					"Age":  &Field{ColName: "age"},
				},
				Columns: map[string]*Field{
					"name": {FieldName: "Name", Typ: reflect.TypeOf("")},
					"age":  {FieldName: "Age", Typ: reflect.TypeOf(18)},
				},
			},
		},
		{
			name: "with 2 pointer, with error",
			entity: func() **Model {
				res := &Model{}
				return &res
			},
			wantErr: errs.ErrPointerOnly,
		},
		{
			name:    "basic type, with error",
			entity:  66,
			wantErr: errs.ErrPointerOnly,
		},
	}

	//构建db
	r := &Registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.ParseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)

		})
	}
}

type testModel struct {
	Name string
	Age  int
}
