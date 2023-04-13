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
		wantRes *model
		wantErr error
	}{
		{
			name:   "normal struct",
			entity: testModel{},
			wantRes: &model{
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
			wantRes: &model{
				tableName: "test_model",
				fields: map[string]*field{
					"Name": &field{colName: "name"},
					"Age":  &field{colName: "age"},
				},
			},
		},
		{
			name: "with 2 pointer, with error",
			entity: func() **model {
				res := &model{}
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
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseModel(tc.entity)
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
