package service

import (
	"errors"

	"goboardapi/internal/config"
	"goboardapi/internal/domain"
	"goboardapi/internal/dto"
	"goboardapi/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("게시글을 찾을 수 없습니다")
)

type PostService struct {
	postRepo repository.PostRepository
	cfg      *config.Config
}

func NewPostService(postRepo repository.PostRepository, cfg *config.Config) *PostService {
	return &PostService{
		postRepo: postRepo,
		cfg:      cfg,
	}
}

func (s *PostService) Create(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	post := &domain.Post{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	return s.toResponse(post), nil
}

func (s *PostService) GetByID(id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	_ = s.postRepo.IncrementViews(id)
	post.Views++

	return s.toResponse(post), nil
}

func (s *PostService) GetList(page, size int, search *dto.SearchParams, sort *dto.SortParams) ([]dto.PostListResponse, *dto.Meta, error) {
	pagination := dto.NewPagination(
		page,
		size,
		s.cfg.Pagination.DefaultSize,
		s.cfg.Pagination.MaxSize,
	)

	posts, total, err := s.postRepo.FindAll(pagination, search, sort)
	if err != nil {
		return nil, nil, err
	}

	list := make([]dto.PostListResponse, len(posts))
	for i, post := range posts {
		list[i] = dto.PostListResponse{
			ID:        post.ID,
			Title:     post.Title,
			Author:    post.Author,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
		}
	}

	totalPages := int(total) / pagination.Size
	if int(total)%pagination.Size > 0 {
		totalPages++
	}

	meta := &dto.Meta{
		Page:       pagination.Page,
		Size:       pagination.Size,
		Total:      total,
		TotalPages: totalPages,
	}

	return list, meta, nil
}

func (s *PostService) GetListByCursor(cursorStr string, size int) ([]dto.PostListResponse, *dto.CursorMeta, error) {
	if size < 1 {
		size = s.cfg.Pagination.DefaultSize
	}
	if size > s.cfg.Pagination.MaxSize {
		size = s.cfg.Pagination.MaxSize
	}

	// 커서 디코딩
	var cursor *dto.Cursor
	if cursorStr != "" {
		var err error
		cursor, err = dto.DecodeCursor(cursorStr)
		if err != nil {
			return nil, nil, errors.New("유효하지 않은 커서입니다")
		}
	}

	// 조회
	posts, err := s.postRepo.FindAllByCursor(cursor, size)
	if err != nil {
		return nil, nil, err
	}

	// 다음 페이지 존재 여부 확인
	hasMore := len(posts) > size
	if hasMore {
		posts = posts[:size]
	}

	// DTO 변환
	list := make([]dto.PostListResponse, len(posts))
	for i, post := range posts {
		list[i] = dto.PostListResponse{
			ID:        post.ID,
			Title:     post.Title,
			Author:    post.Author,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
		}
	}

	// 다음 커서 생성
	var nextCursor string
	if hasMore && len(posts) > 0 {
		last := posts[len(posts)-1]
		c := &dto.Cursor{ID: last.ID, CreatedAt: last.CreatedAt}
		nextCursor = c.Encode()
	}

	meta := &dto.CursorMeta{
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}

	return list, meta, nil
}

func (s *PostService) Update(id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return s.toResponse(post), nil
}

func (s *PostService) Delete(id uint) error {
	_, err := s.postRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	return s.postRepo.Delete(id)
}

func (s *PostService) toResponse(post *domain.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Author:    post.Author,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
