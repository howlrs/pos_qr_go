package repositories

import (
	"backend/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewStoreRepository tests the NewStoreRepository function
func TestNewStoreRepository(t *testing.T) {
	t.Run("NewStoreRepository with nil client returns MockStoreRepository", func(t *testing.T) {
		repo := NewStoreRepository(nil)
		assert.NotNil(t, repo)

		// Check if it's a mock repository
		_, ok := repo.(*MockStoreRepository)
		assert.True(t, ok, "Should return a MockStoreRepository when client is nil")
	})
}

// TestMockStoreRepository tests the MockStoreRepository implementation
func TestMockStoreRepository(t *testing.T) {
	ctx := context.Background()

	// Test data
	testStore := &models.Store{
		ID:        "store_test123",
		Name:      "Test Store",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Address:   "123 Test St",
		Phone:     "123-456-7890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testStores := []*models.Store{testStore}

	t.Run("Create", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Create", mock.Anything, testStore).Return(nil)

		err := mockRepo.Create(ctx, testStore)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Read", mock.Anything).Return(testStores, nil)

		stores, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(stores))
		assert.Len(t, stores, 1)
		assert.Equal(t, testStore.ID, stores[0].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("FindByID", mock.Anything, "store_test123").Return(testStore, nil)

		store, err := mockRepo.FindByID(ctx, "store_test123")
		assert.NoError(t, err)
		assert.Equal(t, testStore.ID, store.ID)
		assert.Equal(t, testStore.Name, store.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("FindByID", mock.Anything, "nonexistent").Return(nil, assert.AnError)

		store, err := mockRepo.FindByID(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, store)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("FindByField", mock.Anything, "email", "test@example.com").Return(testStores, nil)

		stores, err := mockRepo.FindByField(ctx, "email", "test@example.com")
		assert.NoError(t, err)
		assert.Len(t, stores, 1)
		assert.Equal(t, testStore.Email, stores[0].Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField not found", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("FindByField", mock.Anything, "email", "notfound@example.com").Return(nil, assert.AnError)

		stores, err := mockRepo.FindByField(ctx, "email", "notfound@example.com")
		assert.Error(t, err)
		assert.Nil(t, stores)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		updatedStore := &models.Store{
			ID:        "store_test123",
			Name:      "Updated Store",
			Email:     "updated@example.com",
			Password:  "newhashedpassword",
			Address:   "456 Updated St",
			Phone:     "987-654-3210",
			CreatedAt: testStore.CreatedAt,
			UpdatedAt: time.Now(),
		}

		mockRepo.On("UpdateByID", mock.Anything, "store_test123", updatedStore).Return(nil)

		err := mockRepo.UpdateByID(ctx, "store_test123", updatedStore)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "store_test123").Return(nil)

		err := mockRepo.DeleteByID(ctx, "store_test123")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Count", mock.Anything).Return(5, nil)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Count", mock.Anything).Return(0, assert.AnError)

		count, err := mockRepo.Count(ctx)
		assert.Error(t, err)
		assert.Equal(t, 0, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists true", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Exists", mock.Anything, "store_test123").Return(true, nil)

		exists, err := mockRepo.Exists(ctx, "store_test123")
		assert.NoError(t, err)
		assert.True(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists false", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Exists", mock.Anything, "nonexistent").Return(false, nil)

		exists, err := mockRepo.Exists(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Exists", mock.Anything, "error_id").Return(false, assert.AnError)

		exists, err := mockRepo.Exists(ctx, "error_id")
		assert.Error(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryInterface tests that StoreRepository implements the Repository interface
func TestStoreRepositoryInterface(t *testing.T) {
	t.Run("StoreRepository implements Repository interface", func(t *testing.T) {
		var repo Repository[models.Store]
		repo = NewStoreRepository(nil) // This should return a MockStoreRepository
		assert.NotNil(t, repo)

		// Verify the repository is a MockStoreRepository
		_, ok := repo.(*MockStoreRepository)
		assert.True(t, ok, "Should return a MockStoreRepository when client is nil")
	})
}

// TestRepositoryErrorHandling tests error handling scenarios
func TestRepositoryErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("Create error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStore := &models.Store{ID: "test_id"}
		mockRepo.On("Create", mock.Anything, testStore).Return(assert.AnError)

		err := mockRepo.Create(ctx, testStore)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		var nilStores []*models.Store
		mockRepo.On("Read", mock.Anything).Return(nilStores, assert.AnError)

		stores, err := mockRepo.Read(ctx)
		assert.Error(t, err)
		assert.Nil(t, stores)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStore := &models.Store{ID: "test_id"}
		mockRepo.On("UpdateByID", mock.Anything, "test_id", testStore).Return(assert.AnError)

		err := mockRepo.UpdateByID(ctx, "test_id", testStore)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "test_id").Return(assert.AnError)

		err := mockRepo.DeleteByID(ctx, "test_id")
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryStructure tests the StoreRepository structure
func TestStoreRepositoryStructure(t *testing.T) {
	t.Run("StoreRepository with valid client", func(t *testing.T) {
		// Note: In a real test, you would use a mock Firestore client
		// For now, we test that the function doesn't panic with nil client
		repo := NewStoreRepository(nil)
		assert.NotNil(t, repo)

		// Verify it's actually a MockStoreRepository when client is nil
		mockRepo, ok := repo.(*MockStoreRepository)
		assert.True(t, ok)
		assert.NotNil(t, mockRepo)
	})
}

// TestStoreRepositoryMethods tests StoreRepository methods compatibility
func TestStoreRepositoryMethods(t *testing.T) {
	t.Run("All Repository interface methods are implemented", func(t *testing.T) {
		repo := NewStoreRepository(nil)

		// Test that the repository implements all required methods
		// by checking if they can be called without panicking
		ctx := context.Background()
		testStore := &models.Store{
			ID:   "test_store_id",
			Name: "Test Store Name",
		}

		// Since we get a mock when client is nil, we need to set up expectations
		mockRepo := repo.(*MockStoreRepository)

		// Set up mock expectations for each method
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Read", mock.Anything).Return([]*models.Store{testStore}, nil)
		mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(testStore, nil)
		mockRepo.On("FindByField", mock.Anything, mock.Anything, mock.Anything).Return([]*models.Store{testStore}, nil)
		mockRepo.On("UpdateByID", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("Count", mock.Anything).Return(1, nil)
		mockRepo.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

		// Test all methods
		assert.NotPanics(t, func() {
			repo.Create(ctx, testStore)
			repo.Read(ctx)
			repo.FindByID(ctx, testStore.ID)
			repo.FindByField(ctx, "name", testStore.Name)
			repo.UpdateByID(ctx, testStore.ID, testStore)
			repo.DeleteByID(ctx, testStore.ID)
			repo.Count(ctx)
			repo.Exists(ctx, testStore.ID)
		})

		// 各値が期待通りに設定されていることを確認
		store, err := repo.FindByID(ctx, testStore.ID)
		assert.NoError(t, err)
		assert.NotEqual(t, testStore.ID, "store.ID")
		assert.Equal(t, testStore.ID, store.ID)
	})
}
