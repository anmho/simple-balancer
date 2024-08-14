package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	latencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_latency_seconds_total",
			Help:    "Total latency of HTTP requests by method",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method"}, // Label for gRPC methods
	)
)

func metricsMiddleware(handler http.HandlerFunc, handlerName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(w, r)
		latency := time.Since(start)

		requestCounter.WithLabelValues(handlerName).Inc()
		latencyHistogram.WithLabelValues(handlerName).Observe(latency.Seconds())
	}
}

func main() {
	var port int
	if len(os.Args) < 2 {
		log.Fatalln("please provide port to run on")
	}
	x, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln("invalid port value", os.Args[1])
	}
	port = x

	prometheus.MustRegister(latencyHistogram)
	prometheus.MustRegister(requestCounter)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("metrics listening on", port+10)
		log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port+10), nil))
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /",
		metricsMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(
					map[string]any{
						"message": fmt.Sprintf("Hello World from %d", port),
					},
				)
			},
			"HelloWorld",
		),
	)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	log.Printf("starting on port %d\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalln("failed to start server", err)
	}

}
