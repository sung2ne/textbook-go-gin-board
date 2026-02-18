package domain

import (
	"time"

	"gorm.io/gorm"
)

// Comment 댓글 도메인 모델
type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PostID    uint           `gorm:"not null;index" json:"post_id"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`
	AuthorID  uint           `gorm:"not null;index" json:"author_id"`
	Author    *User          `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 연관관계
	Post    Post      `gorm:"foreignKey:PostID" json:"-"`
	Parent  *Comment  `gorm:"foreignKey:ParentID" json:"-"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName 테이블 이름 지정
func (Comment) TableName() string {
	return "comments"
}
