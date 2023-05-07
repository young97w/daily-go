package integration

import (
	"context"
	v1 "geektime/ORM/v14"
	"geektime/ORM/v14/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SelectSuite struct {
	Suite
}

func TestMySQLSelect(t *testing.T) {
	suite.Run(t, &SelectSuite{Suite{
		driver: "mysql",
		dsn:    "root:root@tcp(localhost:13306)/integration_test",
	}})
}

func (s *SelectSuite) TestSelect() {
	t := s.T()
	db := s.db

	testCases := []struct {
		name    string
		builder *v1.Selector[test.SimpleStruct]

		wantErr error
		wantRes *test.SimpleStruct
	}{
		{
			name:    "select id:1",
			builder: v1.NewSelector[test.SimpleStruct](db).Where(v1.C("Id").EQ(1)),
			wantRes: test.NewSimpleStruct(1),
		},
	}

	for _, tc := range testCases {
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
