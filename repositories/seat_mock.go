package repositories

import (
	"backend/models"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockSeatRepository - 実際のFirestoreの複雑な実装は不要
type MockSeatRepository struct {
	mock.Mock
}

func NewMockSeatRepository() Repository[models.Seat] {
	return &MockSeatRepository{}
}

// シンプルな抽象的実装
func (m *MockSeatRepository) Create(ctx context.Context, seat *models.Seat) error {
	args := m.Called(ctx, seat)
	return args.Error(0)
}

func (m *MockSeatRepository) Read(ctx context.Context) ([]*models.Seat, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Seat), args.Error(1)
}

func (m *MockSeatRepository) FindByID(ctx context.Context, id string) (*models.Seat, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Seat), args.Error(1)
}

func (m *MockSeatRepository) FindByField(ctx context.Context, field string, value any) ([]*models.Seat, error) {
	args := m.Called(ctx, field, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Seat), args.Error(1)
}

func (m *MockSeatRepository) UpdateByID(ctx context.Context, id string, seat *models.Seat) error {
	args := m.Called(ctx, id, seat)
	return args.Error(0)
}

func (m *MockSeatRepository) DeleteByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSeatRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int), args.Error(1)
}

func (m *MockSeatRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Get(0).(bool), args.Error(1)
}
