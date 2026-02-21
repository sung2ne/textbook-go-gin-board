package testutil

import (
	"goboardapi/internal/domain"
	"gorm.io/gorm"
)

type Factory struct {
	db *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
	return &Factory{db: db}
}

func (f *Factory) CreateUser(opts ...func(*domain.User)) *domain.User {
	user := &domain.User{
		Username: "기본 사용자",
		Email:    "user@example.com",
	}

	for _, opt := range opts {
		opt(user)
	}

	f.db.Create(user)
	return user
}

func (f *Factory) CreatePost(author *domain.User, opts ...func(*domain.Post)) *domain.Post {
	post := &domain.Post{
		Title:    "기본 제목",
		Content:  "기본 내용",
		AuthorID: author.ID,
	}

	for _, opt := range opts {
		opt(post)
	}

	f.db.Create(post)
	return post
}

func WithUserName(name string) func(*domain.User) {
	return func(u *domain.User) {
		u.Username = name
	}
}

func WithUserEmail(email string) func(*domain.User) {
	return func(u *domain.User) {
		u.Email = email
	}
}

func WithPostTitle(title string) func(*domain.Post) {
	return func(p *domain.Post) {
		p.Title = title
	}
}
