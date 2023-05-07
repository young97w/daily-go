package v1

import (
	"geektime/ORM/internal/errs"
	"geektime/ORM/v14/model"
	"strings"
)

type builder struct {
	sb     strings.Builder
	args   []any
	model  *model.Model
	quoter byte
}

// quote buildColumn addArgs

func (b *builder) quote(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}

func (b *builder) buildColumn(fieldName string) error {

	f, ok := b.model.Fields[fieldName]
	if !ok {
		return errs.NewErrUnknownField(fieldName)
	}
	b.quote(f.ColName)
	return nil
}

func (b *builder) addArgs(args ...any) {
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}

	b.args = append(b.args, args...)
}
