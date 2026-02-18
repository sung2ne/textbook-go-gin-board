package handler

import (
	"net/http"
	"strconv"

	"goboardapi/internal/dto"
	"goboardapi/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// GetNotifications godoc
// @Summary 알림 목록 조회
// @Description 페이징된 알림 목록을 반환합니다
// @Tags notifications
// @Produce json
// @Param page query int false "페이지 번호" default(1)
// @Param size query int false "페이지 크기" default(20)
// @Security Bearer
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response "인증 필요"
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	notifications, total, err := h.notificationService.GetNotifications(
		c.Request.Context(), page, size,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}

	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.SuccessWithMeta(notifications, &dto.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}))
}

// GetUnreadCount godoc
// @Summary 읽지 않은 알림 수 조회
// @Description 읽지 않은 알림 개수를 반환합니다
// @Tags notifications
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.Response "인증 필요"
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	count, err := h.notificationService.GetUnreadCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"unread_count": count,
		},
	})
}

// MarkAsRead godoc
// @Summary 알림 읽음 처리
// @Description 특정 알림을 읽음으로 표시합니다
// @Tags notifications
// @Produce json
// @Param id path int true "알림 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.Response "잘못된 요청"
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID"})
		return
	}

	err = h.notificationService.MarkAsRead(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// MarkAllAsRead godoc
// @Summary 모든 알림 읽음 처리
// @Description 모든 알림을 읽음으로 표시합니다
// @Tags notifications
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} dto.Response "인증 필요"
// @Router /notifications/read-all [put]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	err := h.notificationService.MarkAllAsRead(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
