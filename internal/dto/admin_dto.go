package dto

// AdminStats 관리자 통계
type AdminStats struct {
    TotalUsers    int64 `json:"total_users"`
    TotalPosts    int64 `json:"total_posts"`
    TotalComments int64 `json:"total_comments"`
    TodayUsers    int64 `json:"today_users"`
    TodayPosts    int64 `json:"today_posts"`
}

// ChangeRoleRequest 역할 변경 요청
type ChangeRoleRequest struct {
    Role string `json:"role" binding:"required,oneof=user admin"`
}

// UserListResponse 사용자 목록 응답
type UserListResponse struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}
