package opentelemetry

import (
	"context"
	"fmt"
	v1 "geektime/ORM/v14"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"time"
)

const defaultInstrumentaionName = "geektime/orm/tracing"

type TraceBuilder struct {
	Tracer trace.Tracer
	Name   string

	Shutdown func(context.Context) error
}

func NewTraceBuilder(name string) *TraceBuilder {
	return &TraceBuilder{
		Name: name,
	}
}

func (t *TraceBuilder) WithTracer(tracer trace.Tracer) *TraceBuilder {
	t.Tracer = tracer
	return t
}

func (t *TraceBuilder) Build() v1.Middleware {
	if t.Tracer == nil {
		t.Tracer = otel.GetTracerProvider().Tracer(t.Name)
	}

	return func(next v1.HandleFunc) v1.HandleFunc {
		return func(ctx context.Context, qc *v1.QueryContext) *v1.QueryResult {

			// new span
			spanCtx, span := t.Tracer.Start(ctx, fmt.Sprintf("%s-%s", qc.Type, qc.Model.TableName))

			defer func() {
				t.Shutdown(spanCtx)
			}()
			defer span.End()

			//开始记录信息
			q, _ := qc.Builder.Build()
			if q != nil {
				span.SetAttributes(attribute.String("SQL", q.SQL))
			}

			span.SetAttributes(attribute.String("table", qc.Model.TableName))
			span.SetAttributes(attribute.String("component", "orm"))

			res := next(spanCtx, qc)
			time.Sleep(time.Millisecond * 20)
			if res.Err != nil {
				span.RecordError(res.Err)
			}

			return res
		}
	}
}

func (b *TraceBuilder) WithZipkin(url string) {
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		return
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("orm-test1"),
		)),
	)
	otel.SetTracerProvider(tp)
	b.Shutdown = tp.Shutdown
}

//func (t *TraceBuilder) WithZipkin(url string) *TraceBuilder {
//	exporter, err := zipkin.New(
//		url,
//		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
//	)
//
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	batcher := sdktrace.NewBatchSpanProcessor(exporter)
//	tp := sdktrace.NewTracerProvider(
//		sdktrace.WithSpanProcessor(batcher),
//		sdktrace.WithResource(resource.NewWithAttributes(
//			semconv.SchemaURL,
//			semconv.ServiceNameKey.String(t.Name),
//		)),
//	)
//
//	otel.SetTracerProvider(tp)
//	return t
//}

//WithJaeger 与Jaeger集成
func (t *TraceBuilder) WithJaeger(url string) {
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
				semconv.ServiceNameKey.String(t.Name),
				attribute.String("environment", "dev"),
				attribute.Int64("ID", 1),
			),
		),
	)

	otel.SetTracerProvider(tp)
}
