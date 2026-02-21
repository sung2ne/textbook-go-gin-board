package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func InitDBWithFallback(primaryDSN, replicaDSN string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(primaryDSN), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Replica 연결 시도
    _, replicaErr := gorm.Open(postgres.Open(replicaDSN), &gorm.Config{})

    if replicaErr == nil {
        // Replica 사용 가능
        err = db.Use(dbresolver.Register(dbresolver.Config{
            Replicas: []gorm.Dialector{postgres.Open(replicaDSN)},
            Policy:   dbresolver.RandomPolicy{},
        }))
    } else {
        // Replica 불가 - Primary만 사용
        log.Printf("Replica unavailable, using primary only: %v", replicaErr)
    }

    return db, err
}
