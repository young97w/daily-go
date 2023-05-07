package v1

import (
	"context"
	"fmt"
	"geektime/ORM/v14/internal/test"
	"geektime/ORM/v14/internal/valuer"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRaw(t *testing.T) {
	//构建db
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	//new DB
	db, err := OpenDB(sqlDB)
	db.valCreator = valuer.NewUnsafeValue

	defer sqlDB.Close()
	require.NoError(t, err)

	//构建表数据
	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "young")

	//构建sql表达式返回结果
	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)

	//rows, _ := db.db.QueryContext(context.Background(), "SELECT *")
	//for rows.Next() {
	//	mm := mockModel{}
	//	rows.Scan(&mm.Id, &mm.FirstName)
	//	bytes, _ := json.Marshal(mm)
	//	fmt.Println(string(bytes))
	//}

	testCase := []struct {
		name    string
		builder *RawQuerier[mockModel]

		wantErr error
		wantRes *mockModel
	}{
		{
			name:    "normal model",
			builder: RawQuery[mockModel](db, "SELECT * FROM"),
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

func TestTruncate(t *testing.T) {
	db, err := Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
	require.NoError(t, err)
	q := RawQuery[test.SimpleStruct](db, "TRUNCATE integration_test.simple_struct;")
	res, err := q.Exec(context.Background())
	require.NoError(t, err)
	fmt.Println(res.RowsAffected())
}

func TestRawSelect(t *testing.T) {
	db, err := Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
	require.NoError(t, err)
	q := RawQuery[test.SimpleStruct](db, "SELECT * FROM integration_test.simple_struct;")
	res, err := q.Get(context.Background())
	require.NoError(t, err)
	fmt.Println(res)
}
