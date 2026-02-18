package repository

import (
	"goboardapi/internal/domain"
	"goboardapi/internal/dto"

	"gorm.io/gorm"
)

// PostRepository 게시글 저장소 인터페이스
type PostRepository interface {
	Create(post *domain.Post) error
	FindByID(id uint) (*domain.Post, error)
	FindAll(pagination *dto.Pagination, search *dto.SearchParams, sort *dto.SortParams) ([]domain.Post, int64, error)
	Update(post *domain.Post) error
	Delete(id uint) error
	IncrementViews(id uint) error
	FindAllByCursor(cursor *dto.Cursor, limit int) ([]domain.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll(pagination *dto.Pagination, search *dto.SearchParams, sort *dto.SortParams) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	query := r.db.Model(&domain.Post{})

	// 검색 조건 적용
	if search != nil && search.Query != "" {
		searchQuery := "%" + search.Query + "%"
		switch search.GetSearchType() {
		case dto.SearchTypeTitle:
			query = query.Where("title ILIKE ?", searchQuery)
		case dto.SearchTypeContent:
			query = query.Where("content ILIKE ?", searchQuery)
		default:
			query = query.Where("title ILIKE ? OR content ILIKE ?", searchQuery, searchQuery)
		}
	}

	// 전체 개수 조회
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 정렬 조건 적용
	orderStr := "created_at DESC"
	if sort != nil {
		orderStr = sort.ToOrderString()
	}

	// 페이징 및 정렬 적용하여 조회
	err := query.
		Order(orderStr).
		Offset(pagination.Offset()).
		Limit(pagination.Size).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Post{}, id).Error
}

func (r *postRepository) IncrementViews(id uint) error {
	return r.db.Model(&domain.Post{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).
		Error
}

func (r *postRepository) FindAllByCursor(cursor *dto.Cursor, limit int) ([]domain.Post, error) {
	var posts []domain.Post

	query := r.db.Order("created_at DESC, id DESC")

	// 커서가 있으면 조건 추가
	if cursor != nil {
		query = query.Where(
			"(created_at < ?) OR (created_at = ? AND id < ?)",
			cursor.CreatedAt, cursor.CreatedAt, cursor.ID,
		)
	}

	err := query.Limit(limit + 1).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}
