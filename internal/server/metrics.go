package server

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"go.opentelemetry.io/otel"
	"project/internal/conf"
)

func NewMetricsMiddleware(c *conf.Metrics) middleware.Middleware {
	meter := otel.Meter("project")
	metricRequests, _ := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	metricSeconds, _ := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	return metrics.Server(
		metrics.WithSeconds(metricSeconds),
		metrics.WithRequests(metricRequests),
	)
}