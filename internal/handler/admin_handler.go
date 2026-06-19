package handler

import (
    "errors"
    "net/http"
    "strconv"

    "goboardapi/internal/domain"
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

// GetStats 통계 조회
func (h *AdminHandler) GetStats(c *gin.Context) {
    stats, err := h.adminService.GetStats(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "통계 조회 실패"})
        return
    }

    c.JSON(http.StatusOK, dto.SuccessResponse(stats))
}

// ListUsers 사용자 목록
func (h *AdminHandler) ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

    users, total, err := h.adminService.ListUsers(c.Request.Context(), page, size)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
        return
    }

    var responses []*dto.UserListResponse
    for _, user := range users {
        responses = append(responses, &dto.UserListResponse{
            ID:        user.ID,
            Email:     user.Email,
            Username:  user.Username,
            Role:      string(user.Role),
            CreatedAt: user.CreatedAt,
        })
    }

    c.JSON(http.StatusOK, dto.SuccessWithMeta(responses, &dto.Meta{
        Page:       page,
        Size:       size,
        Total:      total,
        TotalPages: int(total)/size + 1,
    }))
}

// GetUser 사용자 상세
func (h *AdminHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
        return
    }

    user, err := h.adminService.GetUser(c.Request.Context(), uint(id))
    if err != nil {
        if errors.Is(err, repository.ErrUserNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "조회 실패"})
        return
    }

    c.JSON(http.StatusOK, dto.SuccessResponse(user))
}

// ChangeRole 역할 변경
func (h *AdminHandler) ChangeRole(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
        return
    }

    var req dto.ChangeRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newRole := domain.Role(req.Role)
    err = h.adminService.ChangeRole(c.Request.Context(), uint(id), newRole)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "역할이 변경되었습니다",
    })
}

// DeleteUser 사용자 삭제
func (h *AdminHandler) DeleteUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
        return
    }

    err = h.adminService.DeleteUser(c.Request.Context(), uint(id))
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "사용자가 삭제되었습니다",
    })
}

// ForceDeletePost 게시글 강제 삭제
func (h *AdminHandler) ForceDeletePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
        return
    }

    err = h.adminService.ForceDeletePost(c.Request.Context(), uint(id))
    if err != nil {
        if errors.Is(err, repository.ErrPostNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "삭제 실패"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "게시글이 삭제되었습니다",
    })
}

func (h *AdminHandler) handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrCannotDeleteSelf):
        c.JSON(http.StatusBadRequest, gin.H{"error": "자기 자신을 삭제할 수 없습니다"})
    case errors.Is(err, service.ErrCannotChangeSelfRole):
        c.JSON(http.StatusBadRequest, gin.H{"error": "자기 자신의 역할을 변경할 수 없습니다"})
    case errors.Is(err, service.ErrCannotDeleteAdmin):
        c.JSON(http.StatusBadRequest, gin.H{"error": "다른 관리자를 삭제할 수 없습니다"})
    case errors.Is(err, repository.ErrUserNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": "사용자를 찾을 수 없습니다"})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
    }
}
