package opentelemetry

import (
	"context"
	"flag"
	"geektime/ORM/v14"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"testing"
	"time"
)

func TestMiddleware(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "young")

	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)

	//new db
	db, err := v1.OpenDB(mockDB)
	require.NoError(t, err)
	tracer := otel.GetTracerProvider().Tracer("orm-test")
	ot := NewTraceBuilder("orm-test").WithTracer(tracer)
	url := "http://localhost:19411/api/v2/spans"
	ot.WithZipkin(url)

	db.Use(ot.Build())

	testCase := []struct {
		name    string
		builder *v1.Selector[mockModel]

		wantErr error
		wantRes *mockModel
	}{
		{
			name:    "normal model",
			builder: v1.NewSelector[mockModel](db).Where(v1.C("Id").EQ(12)),
			wantRes: &mockModel{
				Id:        1,
				FirstName: "young",
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := context.WithCancel(context.Background())
			res, err := tc.builder.Get(ctx)
			//cancelFunc()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, tc.wantRes, res)
		})
	}
}

type mockModel struct {
	Id        int
	FirstName string
}

func TestTrace(t *testing.T) {
	url := flag.String("zipkin", "http://localhost:19411/api/v2/spans", "zipkin url")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, _ := initTracer(*url)
	defer func() {
		shutdown(ctx)
	}()

	//defer func() {
	//	if err := shutdown(ctx); err != nil {
	//		log.Fatal("failed to shutdown TracerProvider: %w", err)
	//	}
	//}()

	tr := otel.GetTracerProvider().Tracer("component-main")
	ctx, span := tr.Start(ctx, "foo", trace.WithSpanKind(trace.SpanKindServer))
	<-time.After(6 * time.Millisecond)
	bar(ctx)
	<-time.After(6 * time.Millisecond)
	span.End()
}

func bar(ctx context.Context) {
	tr := otel.GetTracerProvider().Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	<-time.After(6 * time.Millisecond)
	span.End()
}

func initTracer(url string) (func(context.Context) error, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(log.New(os.Stderr, "opentelemetry-demo", log.Ldate|log.Ltime|log.Llongfile)),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("zipkin-test345"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
