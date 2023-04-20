package sql_demo

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSqlMock(t *testing.T) {
	//构建db
	db, mock, err := sqlmock.New()
	defer db.Close()
	require.NoError(t, err)

	//构建表数据
	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "young")

	//构建sql表达式返回结果
	mock.ExpectQuery("SELECT `id`,`first_name` FROM `user` .*").WillReturnRows(mockRows)
	mock.ExpectQuery("SELECT `id` FROM `user` .*").WillReturnError(errors.New("mock error"))

	//开始测试
	//第一个expectedQuery
	rows, err := db.QueryContext(context.Background(), "SELECT `id`,`first_name` FROM `user` WHERE `id`=", 1)
	rows.Columns()
	require.NoError(t, err)
	for rows.Next() {
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName)
		require.NoError(t, err)
		fmt.Println(tm)
	}

	//第二个expectedQuery
	_, err = db.QueryContext(context.Background(), "SELECT `id` FROM `user` WHERE `id`=", 1)
	require.Error(t, err)
}
