package v1

import (
	"database/sql"
	"geektime/ORM/v10/internal/model"
	"geektime/ORM/v10/internal/valuer"
)

type DB struct {
	R  *model.Registry
	db *sql.DB

	valCreator valuer.Creator
}

type DBOption func(db *DB)

func Open(driverName, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		R:  &model.Registry{},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}
