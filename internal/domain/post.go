package domain

import (
	"time"

	"gorm.io/gorm"
)

// Post 게시글 도메인 모델
type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	AuthorID  uint           `gorm:"not null;index" json:"author_id"`
	Author    *User          `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Views     int            `gorm:"default:0" json:"views"`
	LikeCount int            `gorm:"default:0" json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 테이블 이름 지정
func (Post) TableName() string {
	return "posts"
}
