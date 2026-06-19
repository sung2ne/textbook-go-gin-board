package repository

type PostRepository interface {
    Create(post *Post) error
    FindByID(id uint) (*Post, error)
    FindAll(page, size int) ([]*Post, error)
    Update(post *Post) error
    Delete(id uint) error
}
