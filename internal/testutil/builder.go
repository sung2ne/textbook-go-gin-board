package testutil

import (
	"goboardapi/internal/domain"
	"gorm.io/gorm"
)

type PostBuilder struct {
	post *domain.Post
	db   *gorm.DB
}

func NewPostBuilder(db *gorm.DB) *PostBuilder {
	return &PostBuilder{
		post: &domain.Post{
			Title:   "기본 제목",
			Content: "기본 내용",
		},
		db: db,
	}
}

func (b *PostBuilder) WithTitle(title string) *PostBuilder {
	b.post.Title = title
	return b
}

func (b *PostBuilder) WithContent(content string) *PostBuilder {
	b.post.Content = content
	return b
}

func (b *PostBuilder) WithAuthor(authorID uint) *PostBuilder {
	b.post.AuthorID = authorID
	return b
}

func (b *PostBuilder) Build() *domain.Post {
	return b.post
}

func (b *PostBuilder) Create() *domain.Post {
	b.db.Create(b.post)
	return b.post
}
