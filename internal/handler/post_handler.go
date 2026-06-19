
// @Summary 게시글 목록 조회
// @Description 페이징된 게시글 목록을 반환합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 크기" default(10)
// @Success 200 {object} dto.ListPostsResponse "성공"
// @Failure 500 {object} dto.ErrorResponse "서버 오류"
// @Router /posts [get]
func (h *PostHandler) List(c *gin.Context) {
    // ...
}

// @Summary 게시글 상세 조회
// @Description ID로 게시글을 조회합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "게시글 ID"
// @Success 200 {object} dto.PostResponse "성공"
// @Failure 400 {object} dto.ErrorResponse "잘못된 ID"
// @Failure 404 {object} dto.ErrorResponse "게시글 없음"
// @Router /posts/{id} [get]
func (h *PostHandler) Get(c *gin.Context) {
    // ...
}

// @Summary 게시글 작성
// @Description 새 게시글을 작성합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "게시글 정보"
// @Success 201 {object} dto.PostResponse "생성됨"
// @Failure 400 {object} dto.ErrorResponse "유효성 검사 실패"
// @Failure 401 {object} dto.ErrorResponse "인증 필요"
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
    // ...
}

// @Summary 게시글 수정
// @Description 게시글을 수정합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "게시글 ID"
// @Param request body dto.UpdatePostRequest true "수정 정보"
// @Success 200 {object} dto.PostResponse "성공"
// @Failure 400 {object} dto.ErrorResponse "유효성 검사 실패"
// @Failure 403 {object} dto.ErrorResponse "권한 없음"
// @Failure 404 {object} dto.ErrorResponse "게시글 없음"
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
    // ...
}

// @Summary 게시글 삭제
// @Description 게시글을 삭제합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "게시글 ID"
// @Success 204 "삭제 완료"
// @Failure 403 {object} dto.ErrorResponse "권한 없음"
// @Failure 404 {object} dto.ErrorResponse "게시글 없음"
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
    // ...
}
