package middleware

import (
    "log"
    "sync"

    "gorm.io/gorm"
)

type QueryCounter struct {
    mu      sync.Mutex
    count   int
    queries []string
}

func (qc *QueryCounter) Add(sql string) {
    qc.mu.Lock()
    defer qc.mu.Unlock()
    qc.count++
    qc.queries = append(qc.queries, sql)
}

func (qc *QueryCounter) Check(threshold int) {
    if qc.count > threshold {
        log.Printf("WARNING: %d queries executed (possible N+1 problem)", qc.count)
        for i, q := range qc.queries {
            log.Printf("  [%d] %s", i+1, q)
        }
    }
}

func SetupNPlusOneDetector(db *gorm.DB) {
    counter := &QueryCounter{}

    db.Callback().Query().After("gorm:query").Register("n_plus_one_check", func(db *gorm.DB) {
        counter.Add(db.Statement.SQL.String())
    })

    // 요청 끝에 확인 (미들웨어에서 호출)
    // counter.Check(10)
}
