
// 관리자 전용 API
admin := r.Group("/api/admin")
admin.Use(middleware.AuthMiddleware(tokenService))
admin.Use(middleware.RequireRole(domain.RoleAdmin))
{
    // 통계
    admin.GET("/stats", adminHandler.GetStats)

    // 사용자 관리
    admin.GET("/users", adminHandler.ListUsers)
    admin.GET("/users/:id", adminHandler.GetUser)
    admin.PUT("/users/:id/role", adminHandler.ChangeRole)
    admin.DELETE("/users/:id", adminHandler.DeleteUser)

    // 콘텐츠 관리
    admin.DELETE("/posts/:id", adminHandler.ForceDeletePost)
    admin.DELETE("/comments/:id", adminHandler.ForceDeleteComment)
}
