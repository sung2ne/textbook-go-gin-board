
// @Summary 내 프로필 조회
// @Description 로그인한 사용자의 프로필을 조회합니다
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {}

// @Summary 프로필 수정
// @Description 로그인한 사용자의 프로필을 수정합니다
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body UpdateProfileRequest true "프로필 정보"
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {}
