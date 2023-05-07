package querylog

import (
	"context"
	"fmt"
	"geektime/ORM/v14"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMiddleware(t *testing.T) {
	fn := func(sql string, args []any) {
		fmt.Println("SQL:", sql)
		fmt.Println("ARGS:", args)
	}

	logFn := NewMiddlewareBuilder().LogFunc(fn).Build()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "young")

	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)

	//new db
	db, err := v1.OpenDB(mockDB)
	require.NoError(t, err)
	db.Use(logFn)

	testCase := []struct {
		name    string
		builder *v1.Selector[mockModel]

		wantErr error
		wantRes *mockModel
	}{
		{
			name:    "normal model",
			builder: v1.NewSelector[mockModel](db).Where(v1.C("Id").EQ(12)),
			wantRes: &mockModel{
				Id:        1,
				FirstName: "young",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.builder.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type mockModel struct {
	Id        int
	FirstName string
}
