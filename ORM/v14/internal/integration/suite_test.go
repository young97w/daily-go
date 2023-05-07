package integration

import (
	v1 "geektime/ORM/v14"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	driver string
	dsn    string
	db     *v1.DB
}

func (s *Suite) SetupSuite() {
	t := s.T()
	db, err := v1.Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
	require.NoError(t, err)
	s.db = db
}
