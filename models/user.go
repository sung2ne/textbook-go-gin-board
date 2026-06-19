package models

import "gorm.io/gorm"

type User struct {
    gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt 포함
    Name  string `gorm:"size:100"`
    Email string `gorm:"uniqueIndex"`
    Age   int
}
