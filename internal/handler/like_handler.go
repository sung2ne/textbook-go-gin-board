package handler

import (
	"errors"
	"net/http"
	"strconv"

	"goboardapi/internal/repository"
	"goboardapi/internal/service"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService service.LikeService
}

func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// LikePost godoc
// @Summary 게시글 좋아요
// @Description 게시글에 좋아요를 누릅니다
// @Tags likes
// @Produce json
// @Param id path int true "게시글 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{} "이미 좋아요"
// @Router /posts/{id}/like [post]
func (h *LikeHandler) LikePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID 형식"})
		return
	}

	err = h.likeService.Like(c.Request.Context(), uint(postID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "좋아요를 눌렀습니다",
	})
}

// UnlikePost godoc
// @Summary 좋아요 취소
// @Description 게시글 좋아요를 취소합니다
// @Tags likes
// @Produce json
// @Param id path int true "게시글 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /posts/{id}/like [delete]
func (h *LikeHandler) UnlikePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID 형식"})
		return
	}

	err = h.likeService.Unlike(c.Request.Context(), uint(postID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "좋아요를 취소했습니다",
	})
}

func (h *LikeHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
	case errors.Is(err, repository.ErrAlreadyLiked):
		c.JSON(http.StatusConflict, gin.H{"error": "이미 좋아요를 눌렀습니다"})
	case errors.Is(err, repository.ErrNotLiked):
		c.JSON(http.StatusConflict, gin.H{"error": "좋아요를 누르지 않았습니다"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}
