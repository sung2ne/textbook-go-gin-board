package middleware

import (
    "log"
    "time"

    "gorm.io/gorm"
)

// QueryLogger - 쿼리 시간 측정 콜백
func SetupQueryLogger(db *gorm.DB) {
    db.Callback().Query().Before("gorm:query").Register("query_start", func(db *gorm.DB) {
        db.InstanceSet("query_start", time.Now())
    })

    db.Callback().Query().After("gorm:query").Register("query_end", func(db *gorm.DB) {
        if start, ok := db.InstanceGet("query_start"); ok {
            elapsed := time.Since(start.(time.Time))
            if elapsed > 100*time.Millisecond {
                log.Printf("[SLOW QUERY] %v: %s", elapsed, db.Statement.SQL.String())
            }
        }
    })
}
