package metrics

import (
    "database/sql"

    "github.com/prometheus/client_golang/prometheus"
)

func RegisterDBMetrics(sqlDB *sql.DB) {
    prometheus.MustRegister(
        prometheus.NewGaugeFunc(
            prometheus.GaugeOpts{
                Name: "db_open_connections",
                Help: "Number of open connections to the database",
            },
            func() float64 {
                return float64(sqlDB.Stats().OpenConnections)
            },
        ),
        prometheus.NewGaugeFunc(
            prometheus.GaugeOpts{
                Name: "db_in_use_connections",
                Help: "Number of connections currently in use",
            },
            func() float64 {
                return float64(sqlDB.Stats().InUse)
            },
        ),
        prometheus.NewGaugeFunc(
            prometheus.GaugeOpts{
                Name: "db_idle_connections",
                Help: "Number of idle connections",
            },
            func() float64 {
                return float64(sqlDB.Stats().Idle)
            },
        ),
    )
}
