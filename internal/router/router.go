
posts := v1.Group("/posts")
{
    posts.GET("", r.postHandler.GetList)
    posts.GET("/cursor", r.postHandler.GetListByCursor)  // 커서 페이징
    posts.POST("", r.postHandler.Create)
    posts.GET("/:id", r.postHandler.GetByID)
    posts.PUT("/:id", r.postHandler.Update)
    posts.DELETE("/:id", r.postHandler.Delete)
}
