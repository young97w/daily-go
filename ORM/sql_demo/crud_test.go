package sql_demo

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestCrud(t *testing.T) {
	//open db
	db, err := sql.Open("sqlite3", "file:mumTest.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	//执行建表
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER ,
    last_name TEXT NOT NULL 
)
`)
	require.NoError(t, err)

	//插入行
	res, err := db.ExecContext(ctx, "INSERT INTO test_model(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)", 1, "young", 18, "sky")
	require.NoError(t, err)

	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响的行数：", affected)

	lastInsertId, err := res.LastInsertId()
	require.NoError(t, err)
	log.Println("最后插入的id：", lastInsertId)

	//查询 并返回单行
	row := db.QueryRowContext(ctx, "SELECT `id`,`first_name`,`age`,`last_name` FROM test_model WHERE `id`=?", 1)
	require.NoError(t, row.Err())
	tm := TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.NoError(t, err)
	fmt.Println(tm)

	//with error
	row = db.QueryRowContext(ctx, "SELECT `id`,`first_name`,`age`,`last_name` FROM test_model WHERE `id`=?", 11)
	require.NoError(t, row.Err())
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.Error(t, sql.ErrNoRows, err)

	//查询多行
	rows, err := db.QueryContext(ctx, "SELECT `id`,`first_name`,`age`,`last_name` FROM test_model WHERE `id`=?", 1)
	require.NoError(t, err)
	//使用迭代器遍历结果
	for rows.Next() {
		tm = TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		fmt.Println(tm)
	}

	cancel()
}

func TestTransaction(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:mumTest.db?cache=shared&mode=memory")
	defer db.Close()
	require.NoError(t, err)

	//set context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	//执行建表
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER ,
    last_name TEXT NOT NULL 
)
`)
	require.NoError(t, err)

	//begin transaction
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)

	//use logic
	res, err := tx.ExecContext(ctx, "INSERT INTO test_model(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?)", 1, "young", 18, "sky")
	require.NoError(t, err)

	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响的行数：", affected)

	lastInsertId, err := res.LastInsertId()
	require.NoError(t, err)
	log.Println("最后插入的id：", lastInsertId)

	//commit
	tx.Commit()

	cancel()
}

func TestPrepareStatement(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:mumTest.db?cache=shared&mode=memory")
	defer db.Close()
	require.NoError(t, err)

	//set context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	//执行建表
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER ,
    last_name TEXT NOT NULL 
)
`)
	require.NoError(t, err)

	//prepare
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` WHERE `id`=?")
	require.NoError(t, err)

	//execute
	rows, err := stmt.QueryContext(ctx, 1)
	require.NoError(t, err)

	//iterate rows
	for rows.Next() {
		tm := TestModel{}
		rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	}

	cancel()
	// 整个应用关闭的时候调用
	stmt.Close()

	// stmt, err = db.PrepareContext(ctx,
	// 	"SELECT * FROM `test_model` WHERE `id` IN (?, ?, ?)")
	// stmt, err = db.PrepareContext(ctx,
	// 	"SELECT * FROM `test_model` WHERE `id` IN (?, ?, ?, ?)")
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
