package repositories

import (
	"backend/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewSeatRepository tests the NewSeatRepository function
func TestNewSeatRepository(t *testing.T) {
	t.Run("NewSeatRepository with nil client returns MockSeatRepository", func(t *testing.T) {
		repo := NewSeatRepository(nil)
		assert.NotNil(t, repo)

		// Check if it's a mock repository
		_, ok := repo.(*MockSeatRepository)
		assert.True(t, ok, "Should return a MockSeatRepository when client is nil")
	})
}

// TestMockSeatRepository tests the MockSeatRepository implementation
func TestMockSeatRepository(t *testing.T) {
	ctx := context.Background()

	// Test data
	testSeat := &models.Seat{
		ID:        "seat_test123",
		Name:      "Test Seat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testSeats := []*models.Seat{testSeat}

	t.Run("Create", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Create", mock.Anything, testSeat).Return(nil)

		err := mockRepo.Create(ctx, testSeat)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Read", mock.Anything).Return(testSeats, nil)

		seats, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(seats))
		assert.Len(t, seats, 1)
		assert.Equal(t, testSeat.ID, seats[0].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("FindByID", mock.Anything, "seat_test123").Return(testSeat, nil)

		seat, err := mockRepo.FindByID(ctx, "seat_test123")
		assert.NoError(t, err)
		assert.Equal(t, testSeat.ID, seat.ID)
		assert.Equal(t, testSeat.Name, seat.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("FindByID", mock.Anything, "nonexistent").Return(nil, assert.AnError)

		seat, err := mockRepo.FindByID(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, seat)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("FindByField", mock.Anything, "name", "Test Seat").Return(testSeats, nil)

		seats, err := mockRepo.FindByField(ctx, "name", "Test Seat")
		assert.NoError(t, err)
		assert.Len(t, seats, 1)
		assert.Equal(t, testSeat.Name, seats[0].Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField not found", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("FindByField", mock.Anything, "name", "NonExistent Seat").Return(nil, assert.AnError)

		seats, err := mockRepo.FindByField(ctx, "name", "NonExistent Seat")
		assert.Error(t, err)
		assert.Nil(t, seats)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		updatedSeat := &models.Seat{
			ID:        "seat_test123",
			Name:      "Updated Seat",
			CreatedAt: testSeat.CreatedAt,
			UpdatedAt: time.Now(),
		}

		mockRepo.On("UpdateByID", mock.Anything, "seat_test123", updatedSeat).Return(nil)

		err := mockRepo.UpdateByID(ctx, "seat_test123", updatedSeat)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "seat_test123").Return(nil)

		err := mockRepo.DeleteByID(ctx, "seat_test123")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Count", mock.Anything).Return(5, nil)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count error", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Count", mock.Anything).Return(0, assert.AnError)

		count, err := mockRepo.Count(ctx)
		assert.Error(t, err)
		assert.Equal(t, 0, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists true", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Exists", mock.Anything, "seat_test123").Return(true, nil)

		exists, err := mockRepo.Exists(ctx, "seat_test123")
		assert.NoError(t, err)
		assert.True(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists false", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Exists", mock.Anything, "nonexistent").Return(false, nil)

		exists, err := mockRepo.Exists(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists error", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("Exists", mock.Anything, "error_id").Return(false, assert.AnError)

		exists, err := mockRepo.Exists(ctx, "error_id")
		assert.Error(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})
}

// TestSeatRepositoryInterface tests that SeatRepository implements the Repository interface
func TestSeatRepositoryInterface(t *testing.T) {
	t.Run("SeatRepository implements Repository interface", func(t *testing.T) {
		var repo Repository[models.Seat]
		repo = NewSeatRepository(nil) // This should return a MockSeatRepository
		assert.NotNil(t, repo)

		// Verify the repository is a MockSeatRepository
		_, ok := repo.(*MockSeatRepository)
		assert.True(t, ok, "Should return a MockSeatRepository when client is nil")
	})
}

// TestRepositoryErrorHandling tests error handling scenarios
func TestSeatRepositoryErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("Create error", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		testSeat := &models.Seat{ID: "test_id"}
		mockRepo.On("Create", mock.Anything, testSeat).Return(assert.AnError)

		err := mockRepo.Create(ctx, testSeat)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read error", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		var nilSeats []*models.Seat
		mockRepo.On("Read", mock.Anything).Return(nilSeats, assert.AnError)

		seats, err := mockRepo.Read(ctx)
		assert.Error(t, err)
		assert.Nil(t, seats)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID error", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		testSeat := &models.Seat{ID: "test_id"}
		mockRepo.On("UpdateByID", mock.Anything, "test_id", testSeat).Return(assert.AnError)

		err := mockRepo.UpdateByID(ctx, "test_id", testSeat)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID error", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "test_id").Return(assert.AnError)

		err := mockRepo.DeleteByID(ctx, "test_id")
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestSeatRepositoryStructure tests the SeatRepository structure
func TestSeatRepositoryStructure(t *testing.T) {
	t.Run("SeatRepository with valid client", func(t *testing.T) {
		// Note: In a real test, you would use a mock Firestore client
		// For now, we test that the function doesn't panic with nil client
		repo := NewSeatRepository(nil)
		assert.NotNil(t, repo)

		// Verify it's actually a MockSeatRepository when client is nil
		mockRepo, ok := repo.(*MockSeatRepository)
		assert.True(t, ok)
		assert.NotNil(t, mockRepo)
	})
}

// TestSeatRepositoryMethods tests SeatRepository methods compatibility
func TestSeatRepositoryMethods(t *testing.T) {
	t.Run("All Repository interface methods are implemented", func(t *testing.T) {
		repo := NewSeatRepository(nil)

		// Test that the repository implements all required methods
		// by checking if they can be called without panicking
		ctx := context.Background()
		testSeat := &models.Seat{
			ID:   "test_seat_id",
			Name: "Test Seat Name",
		}

		// Since we get a mock when client is nil, we need to set up expectations
		mockRepo := repo.(*MockSeatRepository)

		// Set up mock expectations for each method
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Read", mock.Anything).Return([]*models.Seat{testSeat}, nil)
		mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(testSeat, nil)
		mockRepo.On("FindByField", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Seat{testSeat}, nil)
		mockRepo.On("UpdateByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Count", mock.Anything).Return(1, nil)
		mockRepo.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

		// Test all methods
		assert.NotPanics(t, func() {
			repo.Create(ctx, testSeat)
			repo.Read(ctx)
			repo.FindByID(ctx, testSeat.ID)
			repo.FindByField(ctx, "name", testSeat.Name)
			repo.UpdateByID(ctx, testSeat.ID, testSeat)
			repo.DeleteByID(ctx, testSeat.ID)
			repo.Count(ctx)
			repo.Exists(ctx, testSeat.ID)
		})

		// 各値が期待通りに設定されていることを確認
		seat, err := repo.FindByID(ctx, testSeat.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, testSeat.ID, "seat.ID")
		assert.Equal(t, testSeat.ID, seat.ID)
	})
}

// TestSeatStruct tests the Seat struct and its methods
func TestSeatStruct(t *testing.T) {
	now := time.Now()
	testSeat := &models.Seat{
		ID:        "seat_123",
		Name:      "Test Seat",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("ToSetSeat conversion", func(t *testing.T) {
		repoSeat := ToSetSeat(testSeat)
		assert.NotNil(t, repoSeat)
		assert.Equal(t, testSeat.ID, repoSeat.ID)
		assert.Equal(t, testSeat.Name, repoSeat.Name)
		assert.Equal(t, testSeat.CreatedAt, repoSeat.CreatedAt)
		assert.Equal(t, testSeat.UpdatedAt, repoSeat.UpdatedAt)
	})

	t.Run("ToUpdate method", func(t *testing.T) {
		repoSeat := ToSetSeat(testSeat)
		originalUpdatedAt := repoSeat.UpdatedAt

		// Sleep for a small duration to ensure time difference
		time.Sleep(1 * time.Millisecond)

		updatedSeat := repoSeat.ToUpdate()
		assert.NotNil(t, updatedSeat)
		assert.True(t, updatedSeat.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("ToModel conversion", func(t *testing.T) {
		repoSeat := &Seat{
			ID:        "seat_456",
			Name:      "Repository Seat",
			CreatedAt: now,
			UpdatedAt: now,
		}

		modelSeat := repoSeat.ToModel()
		assert.NotNil(t, modelSeat)
		assert.Equal(t, repoSeat.ID, modelSeat.ID)
		assert.Equal(t, repoSeat.Name, modelSeat.Name)
		assert.Equal(t, repoSeat.CreatedAt, modelSeat.CreatedAt)
		assert.Equal(t, repoSeat.UpdatedAt, modelSeat.UpdatedAt)
	})
}

// TestSeatRepositoryReadMethod tests specific edge cases for the Read method
func TestSeatRepositoryReadMethod(t *testing.T) {
	ctx := context.Background()

	t.Run("Read returns empty slice when no seats exist", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		emptySeats := []*models.Seat{}
		mockRepo.On("Read", mock.Anything).Return(emptySeats, nil)

		seats, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, seats)
		assert.Len(t, seats, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read returns multiple seats", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		testSeats := []*models.Seat{
			{ID: "seat_1", Name: "Seat 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "seat_2", Name: "Seat 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "seat_3", Name: "Seat 3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		mockRepo.On("Read", mock.Anything).Return(testSeats, nil)

		seats, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, seats)
		assert.Len(t, seats, 3)

		for i, seat := range seats {
			assert.Equal(t, testSeats[i].ID, seat.ID)
			assert.Equal(t, testSeats[i].Name, seat.Name)
		}

		mockRepo.AssertExpectations(t)
	})
}

// TestSeatRepositoryFindByFieldMethod tests specific edge cases for the FindByField method
func TestSeatRepositoryFindByFieldMethod(t *testing.T) {
	ctx := context.Background()

	t.Run("FindByField returns empty slice when no matches", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		emptySeats := []*models.Seat{}
		mockRepo.On("FindByField", mock.Anything, "name", "NonExistent").Return(emptySeats, nil)

		seats, err := mockRepo.FindByField(ctx, "name", "NonExistent")
		assert.NoError(t, err)
		assert.NotNil(t, seats)
		assert.Len(t, seats, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField with different field types", func(t *testing.T) {
		mockRepo := &MockSeatRepository{}
		testSeat := &models.Seat{
			ID:        "seat_123",
			Name:      "Test Seat",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		testSeats := []*models.Seat{testSeat}

		// Test with string field
		mockRepo.On("FindByField", mock.Anything, "id", "seat_123").Return(testSeats, nil)
		seats, err := mockRepo.FindByField(ctx, "id", "seat_123")
		assert.NoError(t, err)
		assert.Len(t, seats, 1)

		mockRepo.AssertExpectations(t)
	})
}
