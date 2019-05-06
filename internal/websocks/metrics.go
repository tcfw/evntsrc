package websocks

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequest = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "Counter for any http requests",
	}, []string{"stream"})

	socketGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ws_conn",
		Help: "Guage of current websocket connections",
	}, []string{"stream"})

	bytePublishCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_publish_byte_count",
		Help: "Counter of stream event bytes published",
	}, []string{"stream"})

	byteSubscribeCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "event_subscribe_byte_count",
		Help: "Counter of stream event bytes subscribed",
	}, []string{"stream"})
)

func registerMetrics() {
	prometheus.MustRegister(httpRequest)
	prometheus.MustRegister(socketGauge)
	prometheus.MustRegister(bytePublishCounter)
	prometheus.MustRegister(byteSubscribeCounter)
}
