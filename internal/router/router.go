
me := r.Group("/api/me")
me.Use(middleware.AuthMiddleware(tokenService))
{
    me.GET("", userHandler.GetProfile)
    me.PUT("", userHandler.UpdateProfile)
    me.PUT("/password", userHandler.ChangePassword)
    me.DELETE("", userHandler.Withdraw)
    me.GET("/posts", userHandler.GetMyPosts)
    me.GET("/comments", userHandler.GetMyComments)
}
