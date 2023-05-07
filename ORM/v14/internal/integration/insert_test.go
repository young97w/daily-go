package integration

import (
	"context"
	"database/sql"
	v1 "geektime/ORM/v14"
	"geektime/ORM/v14/internal/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InsertSuite struct {
	Suite
}

func TestMySQLInsert(t *testing.T) {
	suite.Run(t, &InsertSuite{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (s *InsertSuite) TestInsert() {
	t := s.T()
	db := s.db

	testCases := []struct {
		name    string
		builder *v1.Inserter[test.SimpleStruct]

		wantErr error
		wantAFR int64
	}{
		{
			name:    "insert one row",
			builder: v1.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(1)),
			wantAFR: 1,
		},
		{
			name:    "insert 3 row",
			builder: v1.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(2), test.NewSimpleStruct(3), test.NewSimpleStruct(4)),
			wantAFR: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.builder.Exec(context.Background())
			assert.Equal(t, tc.wantErr, result.Err)
			if result.Err != nil {
				return
			}
			afr, err := result.Result.(sql.Result).RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, tc.wantAFR, afr)
		})
	}
}
