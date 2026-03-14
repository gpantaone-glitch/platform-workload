package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Simulate occasional errors — 5% error rate
	if rand.Intn(100) < 5 {
		http.Error(w, "internal error", http.StatusInternalServerError)
		requestCount.WithLabelValues(r.Method, r.URL.Path, "500").Inc()
		requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		return
	}

	fmt.Fprintln(w, "ok")
	requestCount.WithLabelValues(r.Method, r.URL.Path, "200").Inc()
	requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "healthy")
	})
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
