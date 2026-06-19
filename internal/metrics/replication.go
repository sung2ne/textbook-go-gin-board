package metrics

import (
    "database/sql"
)

// GetReplicationLag - PostgreSQL 복제 지연 확인
func GetReplicationLag(db *sql.DB) (float64, error) {
    var lag float64
    err := db.QueryRow(`
        SELECT EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp()))
    `).Scan(&lag)
    return lag, err
}
