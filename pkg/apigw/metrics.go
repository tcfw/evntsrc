package apigw

import (
	"net/http"
	"os"

	"github.com/rcrowley/go-metrics"
	influxdb "github.com/vrischmann/go-metrics-influxdb"
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := metrics.GetOrRegisterCounter("httpRequests", metrics.DefaultRegistry)
		m.Inc(1)

		next.ServeHTTP(w, r)
	})
}

func startMetrics() {
	metricsEndpoint, ok := os.LookupEnv("METRICS")
	if !ok {
		metricsEndpoint = "http://metrics:8086"
	}

	tags := map[string]string{
		"app": "apigw",
	}

	if hostname, err := os.Hostname(); err == nil {
		tags["hostname"] = hostname
	}

	go influxdb.InfluxDBWithTags(metrics.DefaultRegistry, 10e9, metricsEndpoint, "evntsrc", "", "", tags)
}
