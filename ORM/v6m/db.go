package v1

import "database/sql"

type DB struct {
	r  *registry
	db *sql.DB
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
		r:  &registry{},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}
