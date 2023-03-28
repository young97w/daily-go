package prometheus

import (
	"geektime/web/v6"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type PrometheusBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func NewPrometheusBuilder(name, subsystem, help string, constLabels map[string]string) *PrometheusBuilder {
	return &PrometheusBuilder{
		Name:        name,
		Subsystem:   subsystem,
		ConstLabels: constLabels,
		Help:        help,
	}
}

func (b *PrometheusBuilder) Build() web.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem:   b.Subsystem,
		Name:        b.Name,
		Help:        b.Help,
		ConstLabels: b.ConstLabels,
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.75: 0.01,
			0.90: 0.01,
			0.99: 0.001,
		},
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summaryVec)
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			startTime := time.Now()
			next(ctx)
			endTime := time.Now()
			go report(endTime.Sub(startTime), ctx, summaryVec)
		}
	}
}

func report(d time.Duration, ctx *web.Context, vec *prometheus.SummaryVec) {
	status := ctx.RespStatusCode
	route := "unknown"
	if ctx.MatchedRoute != "" {
		route = ctx.MatchedRoute
	}
	ms := d / time.Millisecond
	vec.WithLabelValues(route, ctx.Req.Method, strconv.Itoa(status)).Observe(float64(ms))
}
