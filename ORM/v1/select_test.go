package v1

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	testCase := []struct {
		name      string
		builder   *Selector[TestModel]
		wantQuery *Query
		wantErr   error
	}{
		{
			name:      "no from",
			builder:   &Selector[TestModel]{},
			wantQuery: &Query{SQL: "SELECT * FROM TestModel;"},
		},
		{
			name:      "with from",
			builder:   (&Selector[TestModel]{}).From("test_table"),
			wantQuery: &Query{SQL: "SELECT * FROM test_table;"},
		},
		{
			name:      "empty form",
			builder:   NewSelector[TestModel]().From(""),
			wantQuery: &Query{SQL: "SELECT * FROM TestModel;"},
		},
	}

	for _, tc := range testCase {
		query, err := tc.builder.Build()
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			return
		}

		assert.Equal(t, tc.wantQuery, query)
	}
}
