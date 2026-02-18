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

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// Create godoc
// @Summary 댓글 작성
// @Description 게시글에 댓글을 작성합니다
// @Tags comments
// @Accept json
// @Produce json
// @Param postId path int true "게시글 ID"
// @Param request body dto.CreateCommentRequest true "댓글 정보"
// @Security Bearer
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} dto.Response "유효성 검사 실패"
// @Failure 404 {object} dto.Response "게시글 없음"
// @Router /posts/{postId}/comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 게시글 ID입니다"))
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	comment, err := h.commentService.Create(c.Request.Context(), uint(postID), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(dto.ToCommentResponse(comment)))
}

// GetByPostID godoc
// @Summary 댓글 목록 조회
// @Description 게시글의 댓글 목록을 조회합니다
// @Tags comments
// @Produce json
// @Param postId path int true "게시글 ID"
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.Response "게시글 없음"
// @Router /posts/{postId}/comments [get]
func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 게시글 ID입니다"))
		return
	}

	comments, err := h.commentService.GetByPostID(c.Request.Context(), uint(postID))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(comments))
}

// Update godoc
// @Summary 댓글 수정
// @Description 댓글을 수정합니다
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "댓글 ID"
// @Param request body dto.UpdateCommentRequest true "수정 정보"
// @Security Bearer
// @Success 200 {object} dto.CommentResponse
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "댓글 없음"
// @Router /comments/{id} [put]
func (h *CommentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 댓글 ID입니다"))
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	comment, err := h.commentService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToCommentResponse(comment)))
}

// Delete godoc
// @Summary 댓글 삭제
// @Description 댓글을 삭제합니다
// @Tags comments
// @Produce json
// @Param id path int true "댓글 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "댓글 없음"
// @Router /comments/{id} [delete]
func (h *CommentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 댓글 ID입니다"))
		return
	}

	err = h.commentService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "댓글이 삭제되었습니다",
	})
}

func (h *CommentHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "권한이 없습니다"})
	case errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
	case errors.Is(err, repository.ErrCommentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "댓글을 찾을 수 없습니다"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}
