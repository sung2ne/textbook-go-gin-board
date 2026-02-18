package repository

import (
	"context"
	"errors"

	"goboardapi/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]*domain.Notification, int64, error)
	FindUnreadByUserID(ctx context.Context, userID uint) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, id, userID uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	CountUnread(ctx context.Context, userID uint) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]*domain.Notification, int64, error) {
	var notifications []*domain.Notification
	var total int64

	r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ?", userID).
		Count(&total)

	err := r.db.WithContext(ctx).
		Preload("Actor").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) FindUnreadByUserID(ctx context.Context, userID uint) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	err := r.db.WithContext(ctx).
		Preload("Actor").
		Where("user_id = ? AND is_read = false", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)

	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return result.Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}
