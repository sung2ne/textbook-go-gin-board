package models

import (
    "gorm.io/gorm"
)

// User 사용자 모델
type User struct {
    gorm.Model
    Username string `gorm:"size:50;uniqueIndex;not null"`
    Email    string `gorm:"size:100;uniqueIndex;not null"`
    Password string `gorm:"size:255;not null"`
    Role     string `gorm:"size:20;default:'user'"`
    IsActive bool   `gorm:"default:true"`
}

// Post 게시글 모델
type Post struct {
    gorm.Model
    Title    string `gorm:"size:200;not null"`
    Content  string `gorm:"type:text"`
    Views    int    `gorm:"default:0"`
    AuthorID uint   `gorm:"not null;index"`
}

// Comment 댓글 모델
type Comment struct {
    gorm.Model
    Content  string `gorm:"type:text;not null"`
    PostID   uint   `gorm:"not null;index"`
    AuthorID uint   `gorm:"not null;index"`
}
