package repositories

import (
	"backend/models"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockSessionRepository - 実際のFirestoreの複雑な実装は不要
type MockSessionRepository struct {
	mock.Mock
}

func NewMockSessionRepository() Repository[models.Session] {
	return &MockSessionRepository{}
}

// シンプルな抽象的実装
func (m *MockSessionRepository) Create(ctx context.Context, store *models.Session) error {
	args := m.Called(ctx, store)
	return args.Error(0)
}

func (m *MockSessionRepository) Read(ctx context.Context) ([]*models.Session, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return []*models.Session{}, args.Error(1)
	}
	return args.Get(0).([]*models.Session), nil
}

func (m *MockSessionRepository) FindByID(ctx context.Context, id string) (*models.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), nil
}

func (m *MockSessionRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Session, error) {
	args := m.Called(ctx, field, value)
	if args.Get(0) == nil {
		return []*models.Session{}, args.Error(1)
	}
	return args.Get(0).([]*models.Session), nil
}

func (m *MockSessionRepository) UpdateByID(ctx context.Context, id string, store *models.Session) error {
	args := m.Called(ctx, id, store)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), nil
}

func (m *MockSessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), nil
}
