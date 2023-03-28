package errorPage

import "geektime/web/v6"

//创建一个map[int][]byte 来存错误页面或信息
//然后以中间件形式回写数据

type ErrPgBuilder struct {
	m map[int][]byte
}

func NewErrPg() *ErrPgBuilder {
	return &ErrPgBuilder{m: make(map[int][]byte, 64)}
}

func (b *ErrPgBuilder) AddPage(code int, data []byte) {
	b.m[code] = data
}

func (b *ErrPgBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			data, ok := b.m[ctx.RespStatusCode]
			if ok {
				ctx.RespData = data
			}
		}
	}
}
