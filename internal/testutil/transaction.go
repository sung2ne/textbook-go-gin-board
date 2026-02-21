package testutil

import (
    "testing"

    "gorm.io/gorm"
)

// WithTransaction은 트랜잭션 내에서 테스트를 실행하고 롤백합니다.
func WithTransaction(t *testing.T, db *gorm.DB, fn func(tx *gorm.DB)) {
    tx := db.Begin()

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
        tx.Rollback()
    }()

    fn(tx)
}
