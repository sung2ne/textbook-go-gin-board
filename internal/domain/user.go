package domain

import (
	"time"

	"gorm.io/gorm"
)

// Role은 사용자 역할을 나타냅니다.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User는 사용자 엔티티입니다.
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Email       string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Username    string         `gorm:"size:100;not null" json:"username"`
	Role        Role           `gorm:"size:20;default:user" json:"role"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
