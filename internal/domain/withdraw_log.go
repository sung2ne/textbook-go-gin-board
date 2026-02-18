package domain

import "time"

// WithdrawLog 탈퇴 로그
type WithdrawLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Email     string    `gorm:"size:255" json:"email"`
	Reason    string    `gorm:"type:text" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

func (WithdrawLog) TableName() string {
	return "withdraw_logs"
}
