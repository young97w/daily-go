package web

import (
	"encoding/json"
	"fmt"
)

type logFunc func(accessLog string)

type LogFuncBuilder struct {
	logFunc logFunc
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}

func (b *LogFuncBuilder) LogFunc(logFunc logFunc) *LogFuncBuilder {
	b.logFunc = logFunc
	return b
}

//NewLogBuilder 返回结构体
func NewLogBuilder() *LogFuncBuilder {
	return &LogFuncBuilder{logFunc: defaultLogFunc}
}

//Build 返回Middleware （中间件）类型
func (b *LogFuncBuilder) Build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				log := &accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute, //待补充完整
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}

				s, _ := json.Marshal(log)
				b.logFunc(string(s))
			}()
			next(ctx)
		}
	}
}

func defaultLogFunc(s string) {
	fmt.Println(fmt.Sprintf("默认logFunc,%s", s))
}
