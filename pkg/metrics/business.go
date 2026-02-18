package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PostsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "board_posts_total",
		Help: "Total number of posts",
	})

	PostsCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "board_posts_created_total",
		Help: "Total number of posts created",
	})

	CommentsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "board_comments_total",
		Help: "Total number of comments",
	})

	UserLogins = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "board_user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"status"},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "board_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"operation", "table"},
	)
)
