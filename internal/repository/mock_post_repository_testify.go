package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
	"goboardapi/internal/domain"
)

type MockPostRepositoryTestify struct {
	mock.Mock
}

func (m *MockPostRepositoryTestify) Create(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepositoryTestify) FindByID(ctx context.Context, id uint) (*domain.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *MockPostRepositoryTestify) FindAll(ctx context.Context, offset, limit int) ([]*domain.Post, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*domain.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepositoryTestify) Update(ctx context.Context, post *domain.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *MockPostRepositoryTestify) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepositoryTestify) FindByAuthorID(ctx context.Context, authorID uint, offset, limit int) ([]*domain.Post, int64, error) {
	args := m.Called(ctx, authorID, offset, limit)
	return args.Get(0).([]*domain.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepositoryTestify) CountByAuthorID(ctx context.Context, authorID uint) (int64, error) {
	args := m.Called(ctx, authorID)
	return args.Get(0).(int64), args.Error(1)
}
