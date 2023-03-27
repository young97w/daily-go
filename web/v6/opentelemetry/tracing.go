package opentelemetry

import (
	"geektime/web/v6"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"strconv"
)

const defaultInstrumentationName = "geektime/web/v6/tracing"

//以middleware 方式实现tracing
//用户可以自定义tracer

type TraceBuilder struct {
	Name   string
	Tracer trace.Tracer
}

func NewTraceBuilder(name string) *TraceBuilder {
	return &TraceBuilder{Name: name}
}

func (b *TraceBuilder) WithTracer(tracer trace.Tracer) {
	b.Tracer = tracer
}

func (b *TraceBuilder) Build() web.Middleware {
	if b.Tracer == nil {
		b.Tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}

	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			//记录业务函数开始前
			//拿context
			reqCtx := ctx.Req.Context()
			//返回新的context
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))
			//tracing开始
			reqCtx, span := b.Tracer.Start(reqCtx, "unknown", trace.WithAttributes())
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			//关闭span
			defer span.End()

			next(ctx)

			//设置span 名称
			if ctx.MatchedRoute != "" {
				span.SetName(ctx.MatchedRoute)
			}

			//记录响应的http状态码
			span.SetAttributes(attribute.String("http.statucode", strconv.Itoa(ctx.RespStatusCode)))
		}
	}
}

//WithJaeger 与Jaeger集成
func (b *TraceBuilder) WithJaeger(url string) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		log.Fatalln(err)
	}
	tp := sdktrace.NewTracerProvider(
		// Always be sure to batch in production.
		sdktrace.WithBatcher(exporter),
		// Record information about this application in a Resource.
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(b.Name),
				attribute.String("environment", "dev"),
				attribute.Int64("ID", 1),
			),
		),
	)

	otel.SetTracerProvider(tp)
}

//WithZipkin
func (b *TraceBuilder) WithZipkin(url string) {
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		log.Fatalln(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(b.Name),
		)),
	)
	otel.SetTracerProvider(tp)
}
