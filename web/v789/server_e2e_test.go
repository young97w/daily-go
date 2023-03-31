package web

import (
	"fmt"
	"net/http"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Use(mdw1, mdw2, mdw3)
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Get(`/:id(\d+)/:repo`, repo)

	s.Start(":8081")
}

func TestServerWithMiddleware(t *testing.T) {
	serverOpt1 := ServerWithMiddleware(mdw1)
	serverOpt2 := ServerWithMiddleware(mdw2, mdw3)
	//logFunc放最后一个
	l := NewLogBuilder().Build()
	serverOpt3 := ServerWithMiddleware(l)
	s := NewHTTPServer(serverOpt1, serverOpt2, serverOpt3)
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Get(`/:id(\d+)/:repo`, repo)

	err := s.Start(":8081")
	if err != nil {
		panic(err)
	}
}

func repo(ctx *Context) {
	id := ctx.PathParams["id"]
	repoName := ctx.PathParams["repo"]
	res := info{
		Id:   id,
		Repo: repoName,
	}
	ctx.RespJSON(http.StatusOK, res)
}

type info struct {
	Id   string `json:"id"`
	Repo string `json:"repo"`
}

func mdw1(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件1开始")
		next(ctx)
		fmt.Println("中间件1结束")
	}
}

func mdw2(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件2开始")
		next(ctx)
		fmt.Println("中间件2结束")
	}
}

func mdw3(next HandleFunc) HandleFunc {
	return func(ctx *Context) {
		fmt.Println("中间件3开始")
		next(ctx)
		fmt.Println("中间件3结束")
	}
}
