package v1

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func (TestModel) CreateSQL() string {
	return `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`
}

func memoryDB(t *testing.T, opts ...DBOption) *DB {
	db, err := Open("sqlite3",
		"file:test.db?cache=shared&mode=memory",
		// 仅仅用于单元测试，不会发起真的查询
		opts...)
	require.NoError(t, err)
	return db
}
