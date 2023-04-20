package v1

import (
	"geektime/ORM/internal/errs"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				tableName: "test_model",
				fields: map[string]*field{
					"Name": &field{colName: "name"},
					"Age":  &field{colName: "age"},
				},
				columns: map[string]*field{
					"name": {fieldName: "Name", typ: reflect.TypeOf("")},
					"age":  {fieldName: "Age", typ: reflect.TypeOf(18)},
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
				columns: map[string]*field{
					"name": {fieldName: "Name", typ: reflect.TypeOf("")},
					"age":  {fieldName: "Age", typ: reflect.TypeOf(18)},
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
	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	db, _ := OpenDB(sqlDB)
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
