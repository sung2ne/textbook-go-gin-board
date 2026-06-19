package database

import (
    "log"

    "myapp/models"
)

func Migrate() {
    err := DB.AutoMigrate(
        &models.User{},
        &models.Post{},
        &models.Comment{},
    )
    if err != nil {
        log.Fatal("마이그레이션 실패:", err)
    }

    log.Println("마이그레이션 완료")
}
