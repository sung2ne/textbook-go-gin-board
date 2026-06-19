package repository

import (
    "testing"

    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type TransactionalSuite struct {
    suite.Suite
    db *gorm.DB
    tx *gorm.DB
}

func (s *TransactionalSuite) SetupSuite() {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    s.Require().NoError(err)

    err = db.AutoMigrate(&Post{}, &User{}, &Comment{})
    s.Require().NoError(err)

    s.db = db
}

func (s *TransactionalSuite) SetupTest() {
    // 각 테스트 전에 트랜잭션 시작
    s.tx = s.db.Begin()
}

func (s *TransactionalSuite) TearDownTest() {
    // 각 테스트 후에 롤백
    s.tx.Rollback()
}

func (s *TransactionalSuite) TearDownSuite() {
    sqlDB, _ := s.db.DB()
    sqlDB.Close()
}

// DB 접근은 항상 tx를 사용
func (s *TransactionalSuite) DB() *gorm.DB {
    return s.tx
}
