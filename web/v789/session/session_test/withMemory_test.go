package session_test

import (
	"geektime/web/v789"
	"geektime/web/v789/session"
	"geektime/web/v789/session/cookie"
	"geektime/web/v789/session/memery"
	"net/http"
	"testing"
	"time"
)

var m = session.Manager{
	Store:      memery.NewStore(1 * time.Minute),
	Propagator: cookie.NewPropagator("sessid", nil),
	SessCtxKey: "sess",
}

func TestMemory(t *testing.T) {

	s := web.NewHTTPServer()
	s.Use(logMiddleware)

	s.Get("/biz", func(ctx *web.Context) {
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("biz ok")
	})

	s.Post("/logout", func(ctx *web.Context) {
		m.RemoveSession(ctx)
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登出成功")
	})

	s.Post("/login", func(ctx *web.Context) {
		//_, err := m.Generate(ctx.Req.Context(), "young")
		//if err != nil {
		//	ctx.RespStatusCode = http.StatusInternalServerError
		//	ctx.RespData = []byte("系统错误,generate sess")
		//}
		//err = m.Inject("young", ctx.Resp)
		//if err != nil {
		//	ctx.RespStatusCode = http.StatusInternalServerError
		//	ctx.RespData = []byte("系统错误,inject sess")
		//}
		_, err := m.InitSession(ctx, "young")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("系统错误,generate sess")
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功，赋予session")
	})

	s.Start(":8081")

}

func logMiddleware(next web.HandleFunc) web.HandleFunc {
	return func(ctx *web.Context) {
		if ctx.Req.URL.Path != "/login" {
			sess, _ := m.GetSession(ctx)
			//if err != nil {
			//	log.Fatalln(err)
			//}
			if sess == nil {
				ctx.RespStatusCode = http.StatusBadRequest
				ctx.RespData = []byte("请登录系统")
			} else {
				next(ctx)
			}
		} else {
			next(ctx)
		}
	}

}
