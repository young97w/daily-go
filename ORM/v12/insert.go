package v1

import (
	"context"
	"database/sql"
	"geektime/ORM/internal/errs"
	"geektime/ORM/v12/internal/model"
)

type Inserter[T any] struct {
	builder
	db      *DB
	columns []string
	values  []*T
	upsert  *Upsert
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
		builder: builder{
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
	}
}

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

type Upsert struct {
	conflictColumns []string
	assigns         []Assignable
}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

func (u UpsertBuilder[T]) ConflictColumns(cols ...string) UpsertBuilder[T] {
	u.conflictColumns = cols
	return u
}

func (u UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	u.i.upsert = &Upsert{
		conflictColumns: u.conflictColumns,
		assigns:         assigns,
	}
	return u.i
}

func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	i.sb.WriteString("INSERT INTO ")
	m, err := i.db.R.Get(i.values[0])
	if err != nil {
		return nil, err
	}
	i.model = m

	i.quote(m.TableName)

	//build columns
	var cols []*model.Field
	if len(i.columns) > 0 {
		cols = make([]*model.Field, 0, len(i.columns))
		for _, column := range i.columns {
			f, ok := m.Fields[column]
			if !ok {
				return nil, errs.NewErrUnknownField(column)
			}

			cols = append(cols, f)
		}
	} else {
		cols = m.ColSlice
	}

	i.sb.WriteByte('(')
	for c, col := range cols {
		if c > 0 {
			i.sb.WriteByte(',')
		}
		i.quote(col.ColName)
	}
	i.sb.WriteByte(')')

	i.sb.WriteString(" VALUES ")

	//先预估下容量
	i.args = make([]any, 0, len(cols))
	//开始插入值
	for idx, val := range i.values {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		refVal := i.db.valCreator(val, i.model)
		i.sb.WriteByte('(')
		for c, col := range cols {
			if c > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fVal, err := refVal.Field(col.FieldName)
			if err != nil {
				return nil, err
			}

			i.addArgs(fVal)
		}
		i.sb.WriteByte(')')
	}

	//处理 on duplicate key
	if i.upsert != nil {
		i.db.dialect.buildUpsert(&i.builder, i.upsert)
	}
	i.sb.WriteByte(';')
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) Exec(ctx context.Context) (sql.Result, error) {
	q, err := i.Build()
	if err != nil {
		return nil, err
	}

	return i.db.db.ExecContext(ctx, q.SQL, q.Args...)
}
