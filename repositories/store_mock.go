package repositories

import (
	"backend/models"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockStoreRepository - 実際のFirestoreの複雑な実装は不要
type MockStoreRepository struct {
	mock.Mock
}

func NewMockStoreRepository() Repository[models.Store] {
	return &MockStoreRepository{}
}

// シンプルな抽象的実装
func (m *MockStoreRepository) Create(ctx context.Context, store *models.Store) error {
	args := m.Called(ctx, store)
	return args.Error(0)
}

func (m *MockStoreRepository) Read(ctx context.Context) ([]*models.Store, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Store), args.Error(1)
}

func (m *MockStoreRepository) FindByID(ctx context.Context, id string) (*models.Store, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Store), args.Error(1)
}

func (m *MockStoreRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Store, error) {
	args := m.Called(ctx, field, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Store), args.Error(1)
}

func (m *MockStoreRepository) UpdateByID(ctx context.Context, id string, store *models.Store) error {
	args := m.Called(ctx, id, store)
	return args.Error(0)
}

func (m *MockStoreRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStoreRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), args.Error(1)
}

func (m *MockStoreRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), args.Error(1)
}
