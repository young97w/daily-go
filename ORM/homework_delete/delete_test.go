package homework_delete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleter_Build(t *testing.T) {
	testCases := []struct {
		name      string
		builder   *Deleter[TestModel]
		wantErr   error
		wantQuery *Query
	}{
		{
			name:    "no where",
			builder: (&Deleter[TestModel]{}).From("`test_model`"),
			wantQuery: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name:    "where",
			builder: (&Deleter[TestModel]{}).Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `Id` = ?;",
				Args: []any{16},
			},
		},
		{
			name:    "from",
			builder: (&Deleter[TestModel]{}).From("`test_model`").Where(C("Id").EQ(16)),
			wantQuery: &Query{
				SQL:  "DELETE FROM `test_model` WHERE `Id` = ?;",
				Args: []any{16},
			},
		},
	}

	for _, tc := range testCases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			query, err := c.builder.Build()
			assert.Equal(t, c.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
