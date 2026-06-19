package repository

import (
    "gorm.io/plugin/dbresolver"
)

// FindByID - Replica에서 읽기
func (r *PostRepository) FindByID(id uint) (*Post, error) {
    var post Post
    err := r.db.Clauses(dbresolver.Read).First(&post, id).Error
    return &post, err
}

// CreatePost - Primary에 쓰기
func (r *PostRepository) Create(post *Post) error {
    return r.db.Clauses(dbresolver.Write).Create(post).Error
}

// FindByIDFromPrimary - 일관성이 필요할 때 Primary에서 읽기
func (r *PostRepository) FindByIDFromPrimary(id uint) (*Post, error) {
    var post Post
    err := r.db.Clauses(dbresolver.Write).First(&post, id).Error
    return &post, err
}
