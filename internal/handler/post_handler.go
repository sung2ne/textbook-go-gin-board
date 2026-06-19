
// 인증 불필요 - 목록 조회
// @Summary 게시글 목록 조회
// @Tags posts
// @Produce json
// @Success 200 {object} ListPostsResponse
// @Router /posts [get]
func (h *PostHandler) List(c *gin.Context) {}

// 인증 불필요 - 상세 조회
// @Summary 게시글 상세 조회
// @Tags posts
// @Produce json
// @Param id path int true "게시글 ID"
// @Success 200 {object} PostResponse
// @Router /posts/{id} [get]
func (h *PostHandler) Get(c *gin.Context) {}

// 인증 필요 - 작성
// @Summary 게시글 작성
// @Tags posts
// @Accept json
// @Produce json
// @Param request body CreatePostRequest true "게시글 정보"
// @Security Bearer
// @Success 201 {object} PostResponse
// @Failure 401 {object} ErrorResponse "인증 필요"
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {}

// 인증 필요 - 수정
// @Summary 게시글 수정
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "게시글 ID"
// @Param request body UpdatePostRequest true "수정 정보"
// @Security Bearer
// @Success 200 {object} PostResponse
// @Failure 401 {object} ErrorResponse "인증 필요"
// @Failure 403 {object} ErrorResponse "권한 없음"
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {}

// 인증 필요 - 삭제
// @Summary 게시글 삭제
// @Tags posts
// @Produce json
// @Param id path int true "게시글 ID"
// @Security Bearer
// @Success 204 "삭제 완료"
// @Failure 401 {object} ErrorResponse "인증 필요"
// @Failure 403 {object} ErrorResponse "권한 없음"
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {}
