package database

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_USER", "gouser"),
        getEnv("DB_PASSWORD", "gopassword"),
        getEnv("DB_NAME", "godb"),
        getEnv("DB_PORT", "5432"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("데이터베이스 연결 실패:", err)
    }

    sqlDB, _ := DB.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    log.Println("데이터베이스 연결 성공")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
