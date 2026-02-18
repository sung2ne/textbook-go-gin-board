package domain

import (
	"time"
)

// NotificationType 알림 종류
type NotificationType string

const (
	NotificationComment NotificationType = "comment" // 새 댓글
	NotificationReply   NotificationType = "reply"   // 대댓글
	NotificationMention NotificationType = "mention" // 멘션
	NotificationLike    NotificationType = "like"    // 좋아요
)

// Notification 알림 엔티티
type Notification struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    uint             `gorm:"not null;index" json:"user_id"` // 알림 받는 사람
	Type      NotificationType `gorm:"size:20;not null" json:"type"`
	Message   string           `gorm:"size:500;not null" json:"message"`
	Link      string           `gorm:"size:255" json:"link"`
	IsRead    bool             `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`

	// 연관 데이터
	ActorID   uint  `gorm:"index" json:"actor_id"` // 행위자
	PostID    *uint `json:"post_id,omitempty"`
	CommentID *uint `json:"comment_id,omitempty"`

	// 연관관계
	User  User  `gorm:"foreignKey:UserID" json:"-"`
	Actor User  `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
	Post  *Post `gorm:"foreignKey:PostID" json:"-"`
}

// TableName 테이블 이름 지정
func (Notification) TableName() string {
	return "notifications"
}
