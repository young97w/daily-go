package v1

import (
	"geektime/ORM/internal/errs"
	"github.com/stretchr/testify/assert"
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
				tableName: "test_model",
				fields: map[string]*field{
					"Name": &field{colName: "name"},
					"Age":  &field{colName: "age"},
				},
			},
		},
		{
			name:   "with 1 pointer",
			entity: &testModel{},
			wantRes: &Model{
				tableName: "test_model",
				fields: map[string]*field{
					"Name": &field{colName: "name"},
					"Age":  &field{colName: "age"},
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

	db := NewDB()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := db.r.parseModel(tc.entity)
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
