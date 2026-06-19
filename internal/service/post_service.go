package service

func (s *PostService) TransferPost(postID, newUserID uint) error {
    // 트랜잭션 내 모든 쿼리는 Primary
    return s.db.Transaction(func(tx *gorm.DB) error {
        var post model.Post
        if err := tx.First(&post, postID).Error; err != nil {  // Primary에서 읽기
            return err
        }

        post.UserID = newUserID
        return tx.Save(&post).Error  // Primary에 쓰기
    })
}
