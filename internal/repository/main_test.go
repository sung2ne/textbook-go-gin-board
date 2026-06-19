package repository

import (
    "os"
    "testing"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
    // 테스트 전 설정
    var err error
    testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    testDB.AutoMigrate(&Post{}, &Comment{}, &User{})

    // 테스트 실행
    code := m.Run()

    // 테스트 후 정리
    sqlDB, _ := testDB.DB()
    sqlDB.Close()

    os.Exit(code)
}
