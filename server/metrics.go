package main

import (
    "context"
    "net/http"
    "net/http/pprof"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricConnectionsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "c426",
			Name:      "connections_total",
			Help:      "Total number of accepted TCP connections.",
		},
	)

	metricConnectedUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "c426",
			Name:      "connected_users",
			Help:      "Current number of connected authenticated users.",
		},
	)

	metricMessagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "c426",
			Name:      "messages_total",
			Help:      "Total number of messages processed by the server.",
		},
		[]string{"result"}, // success|fail
	)

	metricQueueLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "c426",
			Name:      "queue_length",
			Help:      "Current number of messages in the server queue.",
		},
	)

	metricBlocksIssued = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "c426",
			Name:      "blocks_issued_total",
			Help:      "Total number of blocks issued to users.",
		},
	)

    metricMessageDeliverySeconds = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Namespace: "c426",
            Name:      "message_delivery_seconds",
            Help:      "Time from enqueue to successful delivery.",
            Buckets:   prometheus.DefBuckets,
        },
    )
)

func initMetrics() {
	prometheus.MustRegister(metricConnectionsTotal)
	prometheus.MustRegister(metricConnectedUsers)
	prometheus.MustRegister(metricMessagesSent)
	prometheus.MustRegister(metricQueueLength)
	prometheus.MustRegister(metricBlocksIssued)
    prometheus.MustRegister(metricMessageDeliverySeconds)
}

func startMetricsServer(addr string, shutdownCh <-chan struct{}) {
    srv := &http.Server{
        Addr: addr,
    }
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())
    // pprof endpoints
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
    srv.Handler = mux
    go func() { _ = srv.ListenAndServe() }()
    go func() {
        <-shutdownCh
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()
        _ = srv.Shutdown(ctx)
    }()
}
