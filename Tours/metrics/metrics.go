package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tours_service",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tours_service",
			Name:      "http_request_duration_seconds",
			Help:      "Histogram of request duration.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	// Database connection metrics
	MongoConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "tours_service",
			Name:      "mongo_connections_active",
			Help:      "Number of active MongoDB connections",
		},
		[]string{"database"},
	)

	// Business logic metrics
	ToursCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tours_service",
			Name:      "tours_created_total",
			Help:      "Total number of tours created",
		},
		[]string{"status"},
	)

	KeyPointsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "tours_service",
			Name:      "keypoints_created_total",
			Help:      "Total number of keypoints created",
		},
	)

	ReviewsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "tours_service",
			Name:      "reviews_created_total",
			Help:      "Total number of reviews created",
		},
	)

	TourExecutions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tours_service",
			Name:      "tour_executions_total",
			Help:      "Total number of tour executions",
		},
		[]string{"status"}, // started, finished, abandoned
	)

	ActiveTourExecutions = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "tours_service",
			Name:      "active_tour_executions",
			Help:      "Number of currently active tour executions",
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(MongoConnections)
	prometheus.MustRegister(ToursCreated)
	prometheus.MustRegister(KeyPointsCreated)
	prometheus.MustRegister(ReviewsCreated)
	prometheus.MustRegister(TourExecutions)
	prometheus.MustRegister(ActiveTourExecutions)

	log.Println("Prometheus metrics initialized successfully")
}
