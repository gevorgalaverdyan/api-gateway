package main

import (
	"fmt"
	mockservers "gateway/mock_servers"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "gateway_total_requests",
    Help: "Total number of requests received",
})

var requestsPerRoute = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "gateway_requests_per_route",
    Help: "Requests per route",
}, []string{"route"})

var errorRate = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "gateway_error_rate",
    Help: "Errors per route",
}, []string{"route", "status_code"})

var responseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
    Name:    "gateway_response_time_seconds",
    Help:    "Response time per route",
    Buckets: prometheus.DefBuckets,
}, []string{"route"})

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestsPerRoute)
	prometheus.MustRegister(errorRate)
	prometheus.MustRegister(responseTime)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	routes := map[string]url.URL{
		"/users":    {Scheme: "http", Host: "localhost:9000"},
		"/products": {Scheme: "http", Host: "localhost:9001"},
		"/services": {Scheme: "http", Host: "localhost:9002"},
	}

	hosts := []url.URL{}
	for _, h := range routes {
		hosts = append(hosts, h)
	}

	mockservers.RunServers(hosts)

	
	for rPath, target := range routes {
		http.HandleFunc(rPath, func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// fmt.Fprintf(w, "Forwarding to: %s", target)
			totalRequests.Inc()
			requestsPerRoute.WithLabelValues(r.URL.Path).Inc()
			proxy := httputil.NewSingleHostReverseProxy(&target)
			
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			proxy.ServeHTTP(wrapped, r)
			responseTime.WithLabelValues(r.URL.Path).Observe(time.Since(start).Seconds())

			if wrapped.statusCode >= 400 {
    			errorRate.WithLabelValues(r.URL.Path, fmt.Sprintf("%d", wrapped.statusCode)).Inc()
			}
		})
	}

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
