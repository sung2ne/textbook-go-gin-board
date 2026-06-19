package config

import (
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func InitDB(dsn string) (*gorm.DB, error) {
    // 느린 쿼리 로깅 설정
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             200 * time.Millisecond,  // 200ms 이상 느린 쿼리
            LogLevel:                  logger.Warn,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    return gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
}
