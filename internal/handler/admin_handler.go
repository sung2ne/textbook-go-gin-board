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

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetStats godoc
// @Summary 관리자 통계 조회
// @Description 전체 사이트 통계를 조회합니다
// @Tags admin
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.Response
// @Failure 403 {object} dto.Response "권한 없음"
// @Router /admin/stats [get]
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(stats))
}

// ListUsers godoc
// @Summary 사용자 목록 조회
// @Description 관리자용 사용자 목록을 조회합니다
// @Tags admin
// @Produce json
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 크기" default(20)
// @Security Bearer
// @Success 200 {object} dto.Response
// @Failure 403 {object} dto.Response "권한 없음"
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	users, total, err := h.adminService.ListUsers(c.Request.Context(), page, size)
	if err != nil {
		h.handleError(c, err)
		return
	}

	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(users, &dto.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}))
}

// ChangeRole godoc
// @Summary 사용자 역할 변경
// @Description 사용자의 역할을 변경합니다
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "사용자 ID"
// @Param request body dto.ChangeRoleRequest true "역할 정보"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "사용자 없음"
// @Router /admin/users/{id}/role [put]
func (h *AdminHandler) ChangeRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 사용자 ID"})
		return
	}

	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.adminService.ChangeRole(c.Request.Context(), uint(userID), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "역할이 변경되었습니다",
	})
}

// DeleteUser godoc
// @Summary 사용자 삭제
// @Description 관리자가 사용자를 삭제합니다
// @Tags admin
// @Produce json
// @Param id path int true "사용자 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "사용자 없음"
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 사용자 ID"})
		return
	}

	err = h.adminService.DeleteUser(c.Request.Context(), uint(userID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "사용자가 삭제되었습니다",
	})
}

// ForceDeletePost godoc
// @Summary 게시글 강제 삭제
// @Description 관리자가 게시글을 강제 삭제합니다
// @Tags admin
// @Produce json
// @Param id path int true "게시글 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "게시글 없음"
// @Router /admin/posts/{id} [delete]
func (h *AdminHandler) ForceDeletePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 게시글 ID"})
		return
	}

	err = h.adminService.ForceDeletePost(c.Request.Context(), uint(postID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "게시글이 삭제되었습니다",
	})
}

func (h *AdminHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "권한이 없습니다"})
	case errors.Is(err, repository.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다"})
	case errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}
