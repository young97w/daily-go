package v1

import (
	"context"
	"database/sql"
	"geektime/ORM/internal/errs"
	"geektime/ORM/v14/internal/valuer"
	"geektime/ORM/v14/model"
)

type DB struct {
	db *sql.DB
	core
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
		core: core{
			R:          &model.Registry{},
			dialect:    MySQLDialect,
			valCreator: valuer.NewReflectValue,
		},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx: tx}, nil
}

func (db *DB) DoTx(ctx context.Context, fn func(ctx context.Context, tx *Tx) error, opts *sql.TxOptions) (err error) {
	var tx *Tx
	tx, err = db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		if panicked || err != nil {
			e := tx.Rollback()
			if e != nil {
				err = errs.NewErrFailToRollbackTx(err, e, panicked)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	panicked = false
	return err
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := db.db.ExecContext(ctx, query, args...)
	return res, err
}

func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBWithValCreator(creator valuer.Creator) DBOption {
	return func(db *DB) {
		db.valCreator = creator
	}
}

func (db *DB) Use(mdls ...Middleware) {
	db.mdls = append(db.mdls, mdls...)
}
