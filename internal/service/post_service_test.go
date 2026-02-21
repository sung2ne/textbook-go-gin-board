package service

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "myapp/internal/repository"
)

type MockPostRepo struct {
    mock.Mock
}

func (m *MockPostRepo) Create(post *repository.Post) error {
    args := m.Called(post)
    return args.Error(0)
}

func (m *MockPostRepo) FindByID(id uint) (*repository.Post, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*repository.Post), args.Error(1)
}

// 나머지 메서드...

func TestPostService(t *testing.T) {
    t.Run("Create success", func(t *testing.T) {
        repo := new(MockPostRepo)
        repo.On("Create", mock.AnythingOfType("*repository.Post")).
            Return(nil).
            Run(func(args mock.Arguments) {
                post := args.Get(0).(*repository.Post)
                post.ID = 1
            })

        svc := NewPostService(repo)
        post, err := svc.Create("제목", "내용", 1)

        require.NoError(t, err)
        assert.Equal(t, uint(1), post.ID)
        assert.Equal(t, "제목", post.Title)
        repo.AssertExpectations(t)
    })

    t.Run("Create validation error", func(t *testing.T) {
        repo := new(MockPostRepo)
        // Create가 호출되면 안됨

        svc := NewPostService(repo)
        _, err := svc.Create("", "내용", 1)

        assert.Error(t, err)
        repo.AssertNotCalled(t, "Create")
    })

    t.Run("FindByID success", func(t *testing.T) {
        repo := new(MockPostRepo)
        expected := &repository.Post{ID: 1, Title: "테스트"}
        repo.On("FindByID", uint(1)).Return(expected, nil)

        svc := NewPostService(repo)
        post, err := svc.FindByID(1)

        require.NoError(t, err)
        assert.Equal(t, expected, post)
        repo.AssertExpectations(t)
    })

    t.Run("FindByID not found", func(t *testing.T) {
        repo := new(MockPostRepo)
        repo.On("FindByID", uint(999)).Return(nil, repository.ErrNotFound)

        svc := NewPostService(repo)
        _, err := svc.FindByID(999)

        assert.ErrorIs(t, err, repository.ErrNotFound)
        repo.AssertExpectations(t)
    })
}
