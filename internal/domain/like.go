package domain

import (
	"time"
)

// Like 좋아요 도메인 모델
type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_post" json:"user_id"`
	PostID    uint      `gorm:"not null;uniqueIndex:idx_user_post" json:"post_id"`
	CreatedAt time.Time `json:"created_at"`

	// 연관관계
	User User `gorm:"foreignKey:UserID" json:"-"`
	Post Post `gorm:"foreignKey:PostID" json:"-"`
}

// TableName 테이블 이름 지정
func (Like) TableName() string {
	return "likes"
}
