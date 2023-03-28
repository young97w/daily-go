package prometheus

import (
	"geektime/web/v6"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestPrometheusBuilder_Build(t *testing.T) {
	s := web.NewHTTPServer()
	m := map[string]string{
		"instance_id": "655",
	}
	promeMdl := NewPrometheusBuilder("http_request", "web", "测试", m).Build()

	s.Use(promeMdl)

	s.Get("/user", user)
	//prometheus的观察
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			panic(err)
		}
	}()

	err := s.Start(":8081")
	if err != nil {
		panic(err)
	}
}

func user(ctx *web.Context) {
	n := rand.Intn(1000) + 1
	time.Sleep(time.Duration(n) * time.Millisecond)
	ctx.RespJSON(http.StatusOK, struct {
		Msg string
	}{
		Msg: "hello user",
	})
}
