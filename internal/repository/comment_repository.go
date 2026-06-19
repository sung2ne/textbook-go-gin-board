
// FindByPostIDWithReplies 게시글의 댓글 목록 조회 (대댓글 포함)
func (r *commentRepository) FindByPostIDWithReplies(postID uint) ([]domain.Comment, error) {
    var comments []domain.Comment

    // 최상위 댓글 조회 (대댓글은 Preload로 함께 가져옴)
    err := r.db.
        Where("post_id = ? AND parent_id IS NULL", postID).
        Preload("Replies", func(db *gorm.DB) *gorm.DB {
            return db.Order("created_at ASC")
        }).
        Order("created_at ASC").
        Find(&comments).Error

    if err != nil {
        return nil, err
    }

    return comments, nil
}
