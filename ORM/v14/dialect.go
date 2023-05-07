package v1

import (
	"geektime/ORM/internal/errs"
)

type Dialect interface {
	quoter() byte
	buildUpsert(b *builder, odk *Upsert) error
}

var (
	MySQLDialect   = &mysqlDialect{}
	SQLite3Dialect = &sqlite3Dialect{}
)

type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	return '`'
}

func (s standardSQL) buildUpsert(b *builder, odk *Upsert) error {
	// ON CONFLICT(`COL`) DO UPDATE SET `COL`=
	b.sb.WriteString(" ON CONFLICT")
	if len(odk.conflictColumns) > 0 {
		b.sb.WriteByte('(')
		for i, column := range odk.conflictColumns {
			if i > 0 {
				b.sb.WriteByte(',')
			}

			err := b.buildColumn(column)
			if err != nil {
				return err
			}
		}
		b.sb.WriteByte(')')
	}

	b.sb.WriteString(" DO UPDATE SET ")

	for i, assign := range odk.assigns {
		if i > 0 {
			b.sb.WriteByte(',')
		}

		switch a := assign.(type) {
		case Column:
			err := b.buildColumn(a.name)
			if err != nil {
				return err
			}

			b.sb.WriteString("=excluded.")
			b.buildColumn(a.name)
		case Assignment:
			err := b.buildColumn(a.col)
			if err != nil {
				return err
			}

			b.sb.WriteString("=?")
			b.addArgs(a.val)
		}
	}

	return nil
}

type mysqlDialect struct {
}

func (m *mysqlDialect) quoter() byte {
	return '`'
}

func (m *mysqlDialect) buildUpsert(b *builder, odk *Upsert) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i, assign := range odk.assigns {
		if i > 0 {
			b.sb.WriteByte(',')
		}

		switch a := assign.(type) {
		case Column:
			f, ok := b.model.Fields[a.name]
			if !ok {
				return errs.NewErrUnknownField(a.name)
			}
			b.quote(f.ColName)

			// =values(`col`)
			b.sb.WriteString("=VALUES(`")
			b.sb.WriteString(f.ColName)
			b.sb.WriteString("`)")

		case Assignment:
			err := b.buildColumn(a.col)
			if err != nil {
				return err
			}

			b.sb.WriteString("=?")
			b.addArgs(a.val)
		}

	}
	return nil
}

var _ Dialect = &mysqlDialect{}

// sqlite3Dialect
type sqlite3Dialect struct {
	standardSQL
}

type postgresDialect struct {
	standardSQL
}
