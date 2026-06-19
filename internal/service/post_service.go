package service

func (s *PostService) UpdatePostWriteThrough(ctx context.Context, post *model.Post) error {
    // 1. DB 업데이트
    if err := s.repo.Update(post); err != nil {
        return err
    }

    // 2. 캐시 갱신
    key := fmt.Sprintf("post:%d", post.ID)
    cache.SetRedis(ctx, key, post, 5*time.Minute)

    return nil
}
