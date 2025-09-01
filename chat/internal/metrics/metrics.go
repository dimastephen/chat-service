package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "chatService"
	appName   = "chat"
)

type Metrics struct {
	requestCounter  prometheus.Counter
	responseCounter *prometheus.CounterVec
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_requests_total",
				Help:      "Количество запросов в секунду",
			}),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_response_total",
				Help:      "Количество ответов в секунду",
			}, []string{"service", "status", "method"}),
	}
	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues("auth", status, method).Inc()
}
