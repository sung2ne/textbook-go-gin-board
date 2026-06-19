package handler

type PostHandler struct {
    postService *service.PostService
}

func (h *PostHandler) Create(c *gin.Context) {
    // 1. 요청 파싱
    // 2. 유효성 검사
    // 3. Service 호출
    // 4. 응답 반환
}
