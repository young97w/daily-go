package web

import (
	"net/http"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Get(`/:id(\d+)/:name`, func(c *Context) {
		id := c.PathParams["id"]
		name := c.PathParams["name"]
		res := info{
			Id:   id,
			Repo: name,
		}
		c.RespJSON(http.StatusOK, res)
	})

	s.Start(":8081")
}

//func (c *Context) repo() {
//	id := c.PathParams["id"]
//	res := info{
//		Id:   id,
//		Repo: "",
//	}
//	c.RespJSON(http.StatusOK, res)
//}

type info struct {
	Id   string `json:"id"`
	Repo string `json:"repo"`
}
