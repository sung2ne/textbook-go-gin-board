package repository

import "github.com/stretchr/testify/mock"

type MockPostRepositoryTestify struct {
    mock.Mock
}

func (m *MockPostRepositoryTestify) Create(post *Post) error {
    args := m.Called(post)
    return args.Error(0)
}

func (m *MockPostRepositoryTestify) FindByID(id uint) (*Post, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Post), args.Error(1)
}

func (m *MockPostRepositoryTestify) FindAll(page, size int) ([]*Post, error) {
    args := m.Called(page, size)
    return args.Get(0).([]*Post), args.Error(1)
}

func (m *MockPostRepositoryTestify) Update(post *Post) error {
    args := m.Called(post)
    return args.Error(0)
}

func (m *MockPostRepositoryTestify) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}
