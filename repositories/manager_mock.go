package repositories

import (
	"backend/models"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockManagerRepository - 実際のFirestoreの複雑な実装は不要
type MockManagerRepository struct {
	mock.Mock
}

func NewMockManagerRepository() Repository[models.Manager] {
	return &MockManagerRepository{}
}

// シンプルな抽象的実装
func (m *MockManagerRepository) Create(ctx context.Context, manager *models.Manager) error {
	args := m.Called(ctx, manager)
	return args.Error(0)
}

func (m *MockManagerRepository) Read(ctx context.Context) ([]*models.Manager, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []*models.Manager{}, args.Error(1)
	}
	return args.Get(0).([]*models.Manager), nil
}

func (m *MockManagerRepository) FindByID(ctx context.Context, id string) (*models.Manager, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Manager), nil
}

func (m *MockManagerRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Manager, error) {
	args := m.Called(ctx, field, value)
	if args.Get(0) == nil {
		return []*models.Manager{}, args.Error(1)
	}
	return args.Get(0).([]*models.Manager), nil
}

func (m *MockManagerRepository) UpdateByID(ctx context.Context, id string, manager *models.Manager) error {
	args := m.Called(ctx, id, manager)
	return args.Error(0)
}

func (m *MockManagerRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockManagerRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), nil
}

func (m *MockManagerRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}
