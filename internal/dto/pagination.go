package dto

// Pagination 페이징 요청
type Pagination struct {
	Page int `form:"page" binding:"min=1"`
	Size int `form:"size" binding:"min=1,max=100"`
}

// Offset 오프셋 계산
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

// NewPagination 기본값 적용
func NewPagination(page, size, defaultSize, maxSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = defaultSize
	}
	if size > maxSize {
		size = maxSize
	}
	return &Pagination{Page: page, Size: size}
}
