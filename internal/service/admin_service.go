
func (s *adminService) GetStats(ctx context.Context) (*dto.AdminStats, error) {
    totalUsers, _ := s.userRepo.Count(ctx)
    totalPosts, _ := s.postRepo.Count(ctx)
    totalComments, _ := s.commentRepo.Count(ctx)

    today := time.Now().Truncate(24 * time.Hour)
    todayUsers, _ := s.userRepo.CountSince(ctx, today)
    todayPosts, _ := s.postRepo.CountSince(ctx, today)

    return &dto.AdminStats{
        TotalUsers:    totalUsers,
        TotalPosts:    totalPosts,
        TotalComments: totalComments,
        TodayUsers:    todayUsers,
        TodayPosts:    todayPosts,
    }, nil
}
