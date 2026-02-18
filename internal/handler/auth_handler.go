package handler

import (
	"errors"
	"net/http"
	"strings"

	"goboardapi/internal/auth"
	"goboardapi/internal/dto"
	"goboardapi/internal/middleware"
	"goboardapi/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService  service.AuthService
	tokenService *auth.TokenService
}

func NewAuthHandler(authService service.AuthService, tokenService *auth.TokenService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
	}
}

// Signup godoc
// @Summary 회원가입
// @Description 새 계정을 생성합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "회원가입 정보"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response "유효성 검사 실패"
// @Failure 409 {object} dto.Response "이메일 중복"
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "잘못된 요청 형식입니다",
			"details": err.Error(),
		})
		return
	}

	resp, err := h.authService.Signup(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary 로그인
// @Description 이메일과 비밀번호로 로그인합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "로그인 정보"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response "인증 실패"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청 형식입니다",
		})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "이메일 또는 비밀번호가 올바르지 않습니다",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "서버 오류가 발생했습니다",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary 토큰 갱신
// @Description Refresh Token으로 Access Token을 갱신합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh Token"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response "토큰 만료"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token이 필요합니다",
		})
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrExpiredToken) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "세션이 만료되었습니다. 다시 로그인해주세요",
				"code":  "REFRESH_TOKEN_EXPIRED",
			})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "유효하지 않은 토큰입니다",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary 로그아웃
// @Description 현재 세션을 종료합니다
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.Response "인증 필요"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	claims := middleware.MustGetCurrentUser(c)

	// 현재 Access Token 무효화
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 {
		_ = h.tokenService.RevokeAccessToken(c.Request.Context(), parts[1])
	}

	// Refresh Token 삭제
	_ = h.authService.Logout(c.Request.Context(), claims.UserID)

	c.JSON(http.StatusOK, gin.H{
		"message": "로그아웃되었습니다",
	})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{
			"error": "이미 사용 중인 이메일입니다",
		})

	case errors.Is(err, auth.ErrPasswordTooShort):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "비밀번호는 8자 이상이어야 합니다",
		})

	case errors.Is(err, auth.ErrPasswordNoUpper):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "비밀번호에 대문자가 포함되어야 합니다",
		})

	case errors.Is(err, auth.ErrPasswordNoLower):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "비밀번호에 소문자가 포함되어야 합니다",
		})

	case errors.Is(err, auth.ErrPasswordNoDigit):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "비밀번호에 숫자가 포함되어야 합니다",
		})

	case errors.Is(err, auth.ErrPasswordNoSpecial):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "비밀번호에 특수문자가 포함되어야 합니다",
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "서버 오류가 발생했습니다",
		})
	}
}
