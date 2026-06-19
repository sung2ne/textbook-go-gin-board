package main

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name  string `gorm:"size:100"`
    Email string `gorm:"uniqueIndex"`
    Age   int
}

func main() {
    dsn := "host=localhost user=gouser password=gopassword dbname=godb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&User{})

    // 단건 생성
    user := User{Name: "홍길동", Email: "hong@example.com", Age: 30}
    result := db.Create(&user)

    if result.Error != nil {
        log.Fatal("생성 실패:", result.Error)
    }

    fmt.Printf("생성된 ID: %d\n", user.ID)
    fmt.Printf("영향받은 행: %d\n", result.RowsAffected)
}
