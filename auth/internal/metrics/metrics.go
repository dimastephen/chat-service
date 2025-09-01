package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "chatService"
	appName   = "auth"
)

var metrics *Metrics

type Metrics struct {
	RequestCounter  prometheus.Counter
	ResponseCounter *prometheus.CounterVec
}

func Init(_ context.Context) error {
	metrics = &Metrics{
		RequestCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      appName + "_request_total",
			Subsystem: "grpc",
			Help:      "Количество запросов в секунду",
		}),
		ResponseCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      appName + "_response_total",
			Subsystem: "grpc",
			Help:      "Количество ответов в секунду",
		}, []string{"service", "status", "method"})}

	return nil
}

func IncRequestCounter() {
	metrics.RequestCounter.Inc()
}

func IncResponseCounter(status string, method string) {
	metrics.ResponseCounter.WithLabelValues("auth", status, method).Inc()
}
