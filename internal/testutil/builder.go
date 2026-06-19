package testutil

type PostBuilder struct {
    post *Post
    db   *gorm.DB
}

func NewPostBuilder(db *gorm.DB) *PostBuilder {
    return &PostBuilder{
        post: &Post{
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

func (b *PostBuilder) WithStatus(status string) *PostBuilder {
    b.post.Status = status
    return b
}

func (b *PostBuilder) Published() *PostBuilder {
    b.post.Status = "published"
    return b
}

func (b *PostBuilder) Draft() *PostBuilder {
    b.post.Status = "draft"
    return b
}

func (b *PostBuilder) Build() *Post {
    return b.post
}

func (b *PostBuilder) Create() *Post {
    b.db.Create(b.post)
    return b.post
}
