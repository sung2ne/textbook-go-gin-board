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

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// Create godoc
// @Summary 게시글 작성
// @Description 새 게시글을 작성합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param request body dto.CreatePostRequest true "게시글 정보"
// @Security Bearer
// @Success 201 {object} dto.PostResponse
// @Failure 400 {object} dto.Response "유효성 검사 실패"
// @Failure 401 {object} dto.Response "인증 필요"
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	post, err := h.postService.Create(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(dto.ToPostResponse(post)))
}

// GetByID godoc
// @Summary 게시글 상세 조회
// @Description ID로 게시글을 조회합니다
// @Tags posts
// @Produce json
// @Param id path int true "게시글 ID"
// @Success 200 {object} dto.PostResponse
// @Failure 404 {object} dto.Response "게시글 없음"
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	resp, err := h.postService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(resp))
}

// GetList godoc
// @Summary 게시글 목록 조회
// @Description 페이징된 게시글 목록을 반환합니다
// @Tags posts
// @Produce json
// @Param page query int false "페이지 번호" default(1) minimum(1)
// @Param size query int false "페이지 크기" default(10) minimum(1) maximum(100)
// @Success 200 {object} dto.PostListResponse
// @Router /posts [get]
func (h *PostHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	posts, total, err := h.postService.GetAll(c.Request.Context(), page, size)
	if err != nil {
		h.handleError(c, err)
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

// Update godoc
// @Summary 게시글 수정
// @Description 게시글을 수정합니다
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "게시글 ID"
// @Param request body dto.UpdatePostRequest true "수정 정보"
// @Security Bearer
// @Success 200 {object} dto.PostResponse
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "게시글 없음"
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("VALIDATION_ERROR", err.Error()))
		return
	}

	post, err := h.postService.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToPostResponse(post)))
}

// Delete godoc
// @Summary 게시글 삭제
// @Description 게시글을 삭제합니다
// @Tags posts
// @Produce json
// @Param id path int true "게시글 ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} dto.Response "권한 없음"
// @Failure 404 {object} dto.Response "게시글 없음"
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "유효하지 않은 ID입니다"))
		return
	}

	err = h.postService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "게시글이 삭제되었습니다",
	})
}

func (h *PostHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "인증이 필요합니다"})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "권한이 없습니다"})
	case errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
	}
}
