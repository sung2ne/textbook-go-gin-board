package repository

type MockPostRepository struct {
    posts  map[uint]*Post
    nextID uint
}

func NewMockPostRepository() *MockPostRepository {
    return &MockPostRepository{
        posts:  make(map[uint]*Post),
        nextID: 1,
    }
}

func (m *MockPostRepository) Create(post *Post) error {
    post.ID = m.nextID
    m.nextID++
    m.posts[post.ID] = post
    return nil
}

func (m *MockPostRepository) FindByID(id uint) (*Post, error) {
    post, ok := m.posts[id]
    if !ok {
        return nil, ErrNotFound
    }
    return post, nil
}

// 나머지 메서드 구현...
