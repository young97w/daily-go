package v1

type DB struct {
	r *registry
}

type DBOption func(db *DB)

func NewDB(opts ...DBOption) *DB {
	db := &DB{
		r: &registry{},
	}

	for _, opt := range opts {
		opt(db)
	}

	return db
}
