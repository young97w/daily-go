package opentelemetry

import (
	"geektime/web/v6"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"testing"
	"time"
)

var tracer trace.Tracer

func TestTracing(t *testing.T) {
	//new一个tracer
	tracer = otel.GetTracerProvider().Tracer("tracing-test")
	tb := NewTraceBuilder("tracing-test")
	tb.WithTracer(tracer)
	url := "http://localhost:19411/api/v2/spans"
	tb.WithZipkin(url)

	s := web.NewHTTPServer()

	s.Get("/trace", traceFunc)

	s.Use(tb.Build())
	s.Start(":8081")
}

func traceFunc(ctx *web.Context) {
	c, span := tracer.Start(ctx.Req.Context(), "first_layer")
	defer span.End()

	//子节点二
	c, second := tracer.Start(c, "second_layer")
	time.Sleep(time.Second)
	//子节点3-1
	_, third := tracer.Start(c, "third_1")
	time.Sleep(100 * time.Millisecond)
	third.End()

	//字节点3-2
	_, third2 := tracer.Start(c, "third_1")
	time.Sleep(100 * time.Millisecond)
	third2.End()

	//子节点二计时结束
	second.End()
	ctx.RespJSON(http.StatusOK, struct {
		name string
		id   int
	}{
		name: "hello",
		id:   65535,
	})
}
