package testutil

import "gorm.io/gorm"

// CleanTables는 지정된 테이블들의 데이터를 삭제합니다.
func CleanTables(db *gorm.DB, tables ...string) {
    for _, table := range tables {
        db.Exec("DELETE FROM " + table)
    }
}

// CleanAllTables는 모든 테이블의 데이터를 삭제합니다.
func CleanAllTables(db *gorm.DB) {
    // 외래 키 제약 때문에 순서가 중요
    tables := []string{
        "comments",
        "posts",
        "users",
    }
    CleanTables(db, tables...)
}
