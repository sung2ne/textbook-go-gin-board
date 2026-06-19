package repository

import (
    "context"
    "time"

    "yourproject/internal/domain"
)

func (r *postRepository) FindByID(ctx context.Context, id uint) (*domain.Post, error) {
    // 쿼리 타임아웃 설정
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    var post domain.Post
    err := r.db.WithContext(ctx).First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}
