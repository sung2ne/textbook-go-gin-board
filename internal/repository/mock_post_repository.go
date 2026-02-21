package repository

import (
	"context"

	"goboardapi/internal/domain"
)

type MockPostRepository struct {
	posts  map[uint]*domain.Post
	nextID uint
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[uint]*domain.Post),
		nextID: 1,
	}
}

func (m *MockPostRepository) Create(_ context.Context, post *domain.Post) error {
	post.ID = m.nextID
	m.nextID++
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) FindByID(_ context.Context, id uint) (*domain.Post, error) {
	post, ok := m.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (m *MockPostRepository) FindAll(_ context.Context, offset, limit int) ([]*domain.Post, int64, error) {
	var posts []*domain.Post
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	total := int64(len(posts))
	end := offset + limit
	if end > len(posts) {
		end = len(posts)
	}
	if offset >= len(posts) {
		return []*domain.Post{}, total, nil
	}
	return posts[offset:end], total, nil
}

func (m *MockPostRepository) Update(_ context.Context, post *domain.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(_ context.Context, id uint) error {
	if _, ok := m.posts[id]; !ok {
		return ErrPostNotFound
	}
	delete(m.posts, id)
	return nil
}

func (m *MockPostRepository) FindByAuthorID(_ context.Context, authorID uint, offset, limit int) ([]*domain.Post, int64, error) {
	var posts []*domain.Post
	for _, p := range m.posts {
		if p.AuthorID == authorID {
			posts = append(posts, p)
		}
	}
	total := int64(len(posts))
	end := offset + limit
	if end > len(posts) {
		end = len(posts)
	}
	if offset >= len(posts) {
		return []*domain.Post{}, total, nil
	}
	return posts[offset:end], total, nil
}

func (m *MockPostRepository) CountByAuthorID(_ context.Context, authorID uint) (int64, error) {
	var count int64
	for _, p := range m.posts {
		if p.AuthorID == authorID {
			count++
		}
	}
	return count, nil
}
