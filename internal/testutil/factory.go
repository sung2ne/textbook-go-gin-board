package testutil

import "gorm.io/gorm"

type Factory struct {
    db *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
    return &Factory{db: db}
}

func (f *Factory) CreateUser(opts ...func(*User)) *User {
    user := &User{
        Name:  "기본 사용자",
        Email: "user@example.com",
    }

    for _, opt := range opts {
        opt(user)
    }

    f.db.Create(user)
    return user
}

func (f *Factory) CreatePost(author *User, opts ...func(*Post)) *Post {
    post := &Post{
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

// 옵션 함수들
func WithUserName(name string) func(*User) {
    return func(u *User) {
        u.Name = name
    }
}

func WithUserEmail(email string) func(*User) {
    return func(u *User) {
        u.Email = email
    }
}

func WithPostTitle(title string) func(*Post) {
    return func(p *Post) {
        p.Title = title
    }
}
