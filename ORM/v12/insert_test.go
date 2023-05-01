package v1

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInserter_Exec(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := OpenDB(mockDB)

	testCases := []struct {
		name    string
		builder *Inserter[TestModel]
		wantErr error
		wantAFR int64
	}{
		{
			name: "insert with error",
			builder: func() *Inserter[TestModel] {
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("mock err"))
				return NewInserter[TestModel](db).Columns("FirstName", "Age").Values(
					&TestModel{
						Id:        1,
						FirstName: "young",
						Age:       18,
						LastName: &sql.NullString{
							String: "sky",
							Valid:  true,
						},
					},
				)
			}(),
			wantErr: errors.New("mock err"),
		},
		{
			name: "insert with error",
			builder: func() *Inserter[TestModel] {
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(driver.RowsAffected(1))
				return NewInserter[TestModel](db).Columns("FirstName", "Age").Values(
					&TestModel{
						Id:        1,
						FirstName: "young",
						Age:       18,
						LastName: &sql.NullString{
							String: "sky",
							Valid:  true,
						},
					},
				)
			}(),
			wantAFR: int64(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.builder.Exec(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			rowsAF, err := res.RowsAffected()
			assert.Equal(t, tc.wantAFR, rowsAF)
		})
	}

}

func TestDialectSqlite(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	db, err := OpenDB(mockDB, DBWithDialect(SQLite3Dialect))

	testCases := []struct {
		name      string
		builder   *Inserter[TestModel]
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "sqlite dialect ",
			builder: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					}}).OnDuplicateKey().ConflictColumns("Id").Update(C("Age")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?) ON CONFLICT(`id`) " +
					"DO UPDATE SET `age`=excluded.`age`;",
				Args: []any{int64(1), "young", int8(18), &sql.NullString{String: "sky", Valid: true}},
			},
		},
		{
			name: "sqlite dialect col = val",
			builder: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					}}).OnDuplicateKey().ConflictColumns("Id").Update(Assign("Age", 18)),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?) ON CONFLICT(`id`) " +
					"DO UPDATE SET `age`=?;",
				Args: []any{int64(1), "young", int8(18), &sql.NullString{String: "sky", Valid: true}, 18},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantQuery, q)
		})
	}
}

func TestInserter_Build(t *testing.T) {
	//using default dialect , mysql dialect
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	db, err := OpenDB(mockDB)

	testCases := []struct {
		name      string
		builder   *Inserter[TestModel]
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "insert one record",
			builder: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					},
				},
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					},
				},
			),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?),(?,?,?,?);",
				Args: []any{int64(1), "young", int8(18), &sql.NullString{String: "sky", Valid: true},
					int64(1), "young", int8(18), &sql.NullString{String: "sky", Valid: true}},
			},
		},
		//insert with columns
		{
			name: "with cols",
			builder: NewInserter[TestModel](db).Columns("FirstName", "Age").Values(
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					},
				},
			),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`age`) VALUES (?,?);",
				Args: []any{"young", int8(18)},
			},
		},
		// upsert assign column
		{
			name: "with cols",
			builder: NewInserter[TestModel](db).Columns("FirstName", "Age").Values(
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					},
				},
			).OnDuplicateKey().Update(C("Age")),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`age`) VALUES (?,?) ON DUPLICATE KEY UPDATE `age`=VALUES(`age`);",
				Args: []any{"young", int8(18)},
			},
		},
		// upsert assign value
		{
			name: "with cols",
			builder: NewInserter[TestModel](db).Columns("FirstName", "Age").Values(
				&TestModel{
					Id:        1,
					FirstName: "young",
					Age:       18,
					LastName: &sql.NullString{
						String: "sky",
						Valid:  true,
					},
				},
			).OnDuplicateKey().Update(Assign("Age", 18)),
			wantQuery: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`age`) VALUES (?,?) ON DUPLICATE KEY UPDATE `age`=?;",
				Args: []any{"young", int8(18), 18},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantQuery, q)
		})
	}

}
