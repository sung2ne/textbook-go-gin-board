package handler

import (
	"errors"
	"net/http"
	"strconv"

	"goboardapi/internal/dto"
	"goboardapi/internal/repository"
	"goboardapi/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile 프로필 조회
func (h *UserHandler) GetProfile(c *gin.Context) {
	profile, err := h.userService.GetProfile(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(profile))
}

// UpdateProfile 프로필 수정
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.userService.UpdateProfile(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(profile))
}

// ChangePassword 비밀번호 변경
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), &req)
	if err != nil {
		h.handlePasswordError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "비밀번호가 변경되었습니다",
	})
}

// Withdraw 회원 탈퇴
func (h *UserHandler) Withdraw(c *gin.Context) {
	var req dto.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.Withdraw(c.Request.Context(), &req)
	if err != nil {
		h.handleWithdrawError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "회원 탈퇴가 완료되었습니다. 그동안 이용해 주셔서 감사합니다.",
	})
}

// GetMyPosts 내 게시글 목록
func (h *UserHandler) GetMyPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	posts, total, err := h.userService.GetMyPosts(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}

	var responses []*dto.PostListResponse
	for _, post := range posts {
		responses = append(responses, dto.ToPostListResponse(post))
	}

	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(responses, &dto.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}))
}

// GetMyComments 내 댓글 목록
func (h *UserHandler) GetMyComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	comments, total, err := h.userService.GetMyComments(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}

	var responses []*dto.CommentResponse
	for _, comment := range comments {
		responses = append(responses, dto.ToCommentResponse(comment))
	}

	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(responses, &dto.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}))
}

// SearchUsers 사용자 검색 (멘션 자동완성용)
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusOK, dto.SuccessResponse([]interface{}{}))
		return
	}

	users, err := h.userService.SearchByUsername(c.Request.Context(), query, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
		return
	}

	var results []gin.H
	for _, user := range users {
		results = append(results, gin.H{
			"id":       user.ID,
			"username": user.Username,
		})
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(results))
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrUsernameExists):
		c.JSON(http.StatusConflict, gin.H{"error": "이미 사용 중인 사용자명입니다"})
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}

func (h *UserHandler) handlePasswordError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrWrongPassword):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "현재 비밀번호가 일치하지 않습니다",
			"code":  "WRONG_PASSWORD",
		})
	case errors.Is(err, service.ErrSamePassword):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "새 비밀번호는 현재 비밀번호와 달라야 합니다",
			"code":  "SAME_PASSWORD",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}

func (h *UserHandler) handleWithdrawError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrWrongPassword):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "비밀번호가 일치하지 않습니다",
			"code":  "WRONG_PASSWORD",
		})
	case errors.Is(err, service.ErrCannotWithdrawAdmin):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "관리자는 탈퇴할 수 없습니다. 다른 관리자에게 문의하세요.",
			"code":  "ADMIN_CANNOT_WITHDRAW",
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}
