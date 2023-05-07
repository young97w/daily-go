package v1

import (
	"geektime/ORM/v14/internal/valuer"
	"geektime/ORM/v14/model"
)

type core struct {
	dialect    Dialect
	valCreator valuer.Creator
	R          *model.Registry
	mdls       []Middleware
}
