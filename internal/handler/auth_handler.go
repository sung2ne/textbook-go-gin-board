
// @Summary 로그인
// @Description 이메일과 비밀번호로 로그인합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "로그인 정보"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {}

// @Summary 회원가입
// @Description 새 계정을 생성합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SignupRequest true "회원가입 정보"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {}

// @Summary 토큰 갱신
// @Description Refresh Token으로 Access Token을 갱신합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh Token"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {}
