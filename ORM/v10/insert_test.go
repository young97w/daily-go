package v1

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInserter_Build(t *testing.T) {
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
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES (?,?,?,?) (?,?,?,?);",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantQuery.Args, q.Args)
		})
	}

}

func TestTestModel(t *testing.T) {
	t1 := TestModel{LastName: &sql.NullString{
		String: "young",
		Valid:  true,
	}}

	assert.Equal(t, &sql.NullString{
		String: "young",
		Valid:  true,
	}, t1.LastName)
}
