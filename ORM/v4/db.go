package v1

import "reflect"

type DB struct {
	r *registry
}

type DBOption func(db *DB)

func NewRegistry(opts ...DBOption) *DB {
	db := &DB{
		r: &registry{models: make(map[reflect.Type]*model, 64)},
	}

	for _, opt := range opts {
		opt(db)
	}

	return db
}
