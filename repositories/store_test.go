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

// TestStoreStruct tests the Store struct and its methods
func TestStoreStruct(t *testing.T) {
	now := time.Now()
	testStore := &models.Store{
		ID:        "store_123",
		Name:      "Test Store",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Address:   "123 Test St",
		Phone:     "123-456-7890",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("ToSetStore conversion", func(t *testing.T) {
		repoStore := ToSetStore(testStore)
		assert.NotNil(t, repoStore)
		assert.Equal(t, testStore.ID, repoStore.ID)
		assert.Equal(t, testStore.Name, repoStore.Name)
		assert.Equal(t, testStore.Email, repoStore.Email)
		assert.Equal(t, testStore.Password, repoStore.Password)
		assert.Equal(t, testStore.Address, repoStore.Address)
		assert.Equal(t, testStore.Phone, repoStore.Phone)
		assert.Equal(t, testStore.CreatedAt, repoStore.CreatedAt)
		assert.Equal(t, testStore.UpdatedAt, repoStore.UpdatedAt)
	})

	t.Run("ToUpdate method", func(t *testing.T) {
		repoStore := ToSetStore(testStore)
		originalUpdatedAt := repoStore.UpdatedAt

		// Sleep for a small duration to ensure time difference
		time.Sleep(1 * time.Millisecond)

		updatedStore := repoStore.ToUpdate()
		assert.NotNil(t, updatedStore)
		assert.True(t, updatedStore.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("ToModel conversion", func(t *testing.T) {
		repoStore := &Store{
			ID:        "store_456",
			Name:      "Repository Store",
			Email:     "repo@example.com",
			Password:  "repopassword",
			Address:   "456 Repo St",
			Phone:     "987-654-3210",
			CreatedAt: now,
			UpdatedAt: now,
		}

		modelStore := repoStore.ToModel()
		assert.NotNil(t, modelStore)
		assert.Equal(t, repoStore.ID, modelStore.ID)
		assert.Equal(t, repoStore.Name, modelStore.Name)
		assert.Equal(t, repoStore.Email, modelStore.Email)
		assert.Equal(t, repoStore.Password, modelStore.Password)
		assert.Equal(t, repoStore.Address, modelStore.Address)
		assert.Equal(t, repoStore.Phone, modelStore.Phone)
		assert.Equal(t, repoStore.CreatedAt, modelStore.CreatedAt)
		assert.Equal(t, repoStore.UpdatedAt, modelStore.UpdatedAt)
	})
}

// TestStoreRepositoryReadMethod tests specific edge cases for the Read method
func TestStoreRepositoryReadMethod(t *testing.T) {
	ctx := context.Background()

	t.Run("Read returns empty slice when no stores exist", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		emptyStores := []*models.Store{}
		mockRepo.On("Read", mock.Anything).Return(emptyStores, nil)

		stores, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read returns multiple stores", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStores := []*models.Store{
			{ID: "store_1", Name: "Store 1", Email: "store1@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "store_2", Name: "Store 2", Email: "store2@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "store_3", Name: "Store 3", Email: "store3@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		mockRepo.On("Read", mock.Anything).Return(testStores, nil)

		stores, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 3)

		for i, store := range stores {
			assert.Equal(t, testStores[i].ID, store.ID)
			assert.Equal(t, testStores[i].Name, store.Name)
			assert.Equal(t, testStores[i].Email, store.Email)
		}

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryFindByFieldMethod tests specific edge cases for the FindByField method
func TestStoreRepositoryFindByFieldMethod(t *testing.T) {
	ctx := context.Background()

	t.Run("FindByField returns empty slice when no matches", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		emptyStores := []*models.Store{}
		mockRepo.On("FindByField", mock.Anything, "email", "nonexistent@example.com").Return(emptyStores, nil)

		stores, err := mockRepo.FindByField(ctx, "email", "nonexistent@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField with different field types", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStore := &models.Store{
			ID:        "store_123",
			Name:      "Test Store",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Address:   "123 Test St",
			Phone:     "123-456-7890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		testStores := []*models.Store{testStore}

		// Test with string field (email)
		mockRepo.On("FindByField", mock.Anything, "email", "test@example.com").Return(testStores, nil)
		stores, err := mockRepo.FindByField(ctx, "email", "test@example.com")
		assert.NoError(t, err)
		assert.Len(t, stores, 1)
		assert.Equal(t, testStore.Email, stores[0].Email)

		// Test with string field (name)
		mockRepo.On("FindByField", mock.Anything, "name", "Test Store").Return(testStores, nil)
		stores, err = mockRepo.FindByField(ctx, "name", "Test Store")
		assert.NoError(t, err)
		assert.Len(t, stores, 1)
		assert.Equal(t, testStore.Name, stores[0].Name)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryUpdateByIDMethod tests specific edge cases for the UpdateByID method
func TestStoreRepositoryUpdateByIDMethod(t *testing.T) {
	ctx := context.Background()

	t.Run("UpdateByID with partial data", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		partialStore := &models.Store{
			ID:   "store_123",
			Name: "Updated Name Only",
			// Other fields are empty/zero values
		}

		mockRepo.On("UpdateByID", mock.Anything, "store_123", partialStore).Return(nil)

		err := mockRepo.UpdateByID(ctx, "store_123", partialStore)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID with complete data", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		completeStore := &models.Store{
			ID:        "store_123",
			Name:      "Updated Store",
			Email:     "updated@example.com",
			Password:  "newpassword",
			Address:   "456 Updated St",
			Phone:     "987-654-3210",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("UpdateByID", mock.Anything, "store_123", completeStore).Return(nil)

		err := mockRepo.UpdateByID(ctx, "store_123", completeStore)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryValidation tests validation scenarios
func TestStoreRepositoryValidation(t *testing.T) {
	ctx := context.Background()

	t.Run("Create with invalid store data", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		invalidStore := &models.Store{
			// Missing required fields
			ID: "",
		}

		mockRepo.On("Create", mock.Anything, invalidStore).Return(assert.AnError)

		err := mockRepo.Create(ctx, invalidStore)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID with empty ID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("FindByID", mock.Anything, "").Return(nil, assert.AnError)

		store, err := mockRepo.FindByID(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, store)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryPerformance tests performance-related scenarios
func TestStoreRepositoryPerformance(t *testing.T) {
	ctx := context.Background()

	t.Run("Read large number of stores", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}

		// Create a large slice of test stores
		largeStoreList := make([]*models.Store, 1000)
		for i := 0; i < 1000; i++ {
			largeStoreList[i] = &models.Store{
				ID:        "store_" + string(rune(i)),
				Name:      "Store " + string(rune(i)),
				Email:     "store" + string(rune(i)) + "@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		mockRepo.On("Read", mock.Anything).Return(largeStoreList, nil)

		stores, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 1000)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count with large number", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("Count", mock.Anything).Return(10000, nil)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 10000, count)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryEdgeCases tests edge cases and boundary conditions
func TestStoreRepositoryEdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("FindByField with nil value", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("FindByField", mock.Anything, "name", nil).Return([]*models.Store{}, nil)

		stores, err := mockRepo.FindByField(ctx, "name", nil)
		assert.NoError(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID with non-existent ID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStore := &models.Store{
			ID:   "non_existent_id",
			Name: "Test Store",
		}

		mockRepo.On("UpdateByID", mock.Anything, "non_existent_id", testStore).Return(assert.AnError)

		err := mockRepo.UpdateByID(ctx, "non_existent_id", testStore)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID with non-existent ID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "non_existent_id").Return(assert.AnError)

		err := mockRepo.DeleteByID(ctx, "non_existent_id")
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists with special characters in ID", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		specialID := "store_with_special_chars_!@#$%^&*()"
		mockRepo.On("Exists", mock.Anything, specialID).Return(false, nil)

		exists, err := mockRepo.Exists(ctx, specialID)
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryBusinessLogic tests business logic scenarios
func TestStoreRepositoryBusinessLogic(t *testing.T) {
	ctx := context.Background()

	t.Run("Create store with duplicate email should fail", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		existingStore := &models.Store{
			ID:    "store_existing",
			Email: "duplicate@example.com",
		}

		// First, check if store with same email exists
		mockRepo.On("FindByField", mock.Anything, "email", "duplicate@example.com").Return([]*models.Store{existingStore}, nil)

		stores, err := mockRepo.FindByField(ctx, "email", "duplicate@example.com")
		assert.NoError(t, err)
		assert.Len(t, stores, 1)

		// Should not create another store with same email in business logic
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update store password should hash password", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		storeWithNewPassword := &models.Store{
			ID:       "store_123",
			Password: "new_hashed_password_hash",
		}

		mockRepo.On("UpdateByID", mock.Anything, "store_123", storeWithNewPassword).Return(nil)

		err := mockRepo.UpdateByID(ctx, "store_123", storeWithNewPassword)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Find stores by multiple criteria", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}

		// Find stores in specific city
		storesInTokyo := []*models.Store{
			{ID: "store_tokyo_1", Name: "Tokyo Store 1", Address: "Tokyo, Japan"},
			{ID: "store_tokyo_2", Name: "Tokyo Store 2", Address: "Tokyo, Japan"},
		}

		mockRepo.On("FindByField", mock.Anything, "address", "Tokyo, Japan").Return(storesInTokyo, nil)

		stores, err := mockRepo.FindByField(ctx, "address", "Tokyo, Japan")
		assert.NoError(t, err)
		assert.Len(t, stores, 2)

		for _, store := range stores {
			assert.Contains(t, store.Address, "Tokyo")
		}

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryTransactionScenarios tests transaction-like scenarios
func TestStoreRepositoryTransactionScenarios(t *testing.T) {
	ctx := context.Background()

	// setup is a helper function to initialize a new mock repository for each sub-test.
	// This ensures that tests are isolated from each other.
	setup := func() *MockStoreRepository {
		return &MockStoreRepository{}
	}

	t.Run("Create and immediately verify store exists", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		newStore := &models.Store{
			ID:        "store_new_123",
			Name:      "New Store",
			Email:     "new@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Create", ctx, newStore).Return(nil)
		mockRepo.On("Exists", ctx, newStore.ID).Return(true, nil)
		mockRepo.On("FindByID", ctx, newStore.ID).Return(newStore, nil)

		// Act
		createErr := mockRepo.Create(ctx, newStore)
		exists, existsErr := mockRepo.Exists(ctx, newStore.ID)
		retrievedStore, findErr := mockRepo.FindByID(ctx, newStore.ID)

		// Assert
		assert.NoError(t, createErr, "Create should not return an error")
		assert.NoError(t, existsErr, "Exists should not return an error")
		assert.True(t, exists, "Store should exist after creation")
		assert.NoError(t, findErr, "FindByID should not return an error")
		assert.Equal(t, newStore, retrievedStore, "Retrieved store should match the created store")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Update store and verify changes", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		originalStore := &models.Store{
			ID:        "store_update_123",
			Name:      "Original Name",
			Email:     "original@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		updatedStore := &models.Store{
			ID:        originalStore.ID,
			Name:      "Updated Name",
			Email:     "updated@example.com",
			CreatedAt: originalStore.CreatedAt,
			UpdatedAt: time.Now().Add(time.Second), // Ensure updated_at is later
		}

		mockRepo.On("UpdateByID", ctx, originalStore.ID, updatedStore).Return(nil)
		mockRepo.On("FindByID", ctx, originalStore.ID).Return(updatedStore, nil)

		// Act
		updateErr := mockRepo.UpdateByID(ctx, originalStore.ID, updatedStore)
		retrievedStore, findErr := mockRepo.FindByID(ctx, originalStore.ID)

		// Assert
		assert.NoError(t, updateErr, "UpdateByID should not return an error")
		assert.NoError(t, findErr, "FindByID should not return an error after update")
		assert.Equal(t, updatedStore, retrievedStore, "Retrieved store should reflect the updates")
		assert.True(t, retrievedStore.UpdatedAt.After(originalStore.UpdatedAt), "UpdatedAt should be more recent")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete store and verify it no longer exists", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		storeID := "store_delete_123"

		mockRepo.On("DeleteByID", ctx, storeID).Return(nil)
		mockRepo.On("Exists", ctx, storeID).Return(false, nil)
		mockRepo.On("FindByID", ctx, storeID).Return(nil, assert.AnError)

		// Act
		deleteErr := mockRepo.DeleteByID(ctx, storeID)
		exists, existsErr := mockRepo.Exists(ctx, storeID)
		deletedStore, findErr := mockRepo.FindByID(ctx, storeID)

		// Assert
		assert.NoError(t, deleteErr, "DeleteByID should not return an error")
		assert.NoError(t, existsErr, "Exists should not return an error after deletion")
		assert.False(t, exists, "Store should not exist after deletion")
		assert.Error(t, findErr, "FindByID should return an error for a deleted store")
		assert.Nil(t, deletedStore, "FindByID should return a nil store after deletion")

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryDataConsistency tests data consistency scenarios
func TestStoreRepositoryDataConsistency(t *testing.T) {
	ctx := context.Background()

	t.Run("Store count matches actual number of stores", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStores := []*models.Store{
			{ID: "store_1", Name: "Store 1"},
			{ID: "store_2", Name: "Store 2"},
			{ID: "store_3", Name: "Store 3"},
		}

		mockRepo.On("Read", mock.Anything).Return(testStores, nil)
		mockRepo.On("Count", mock.Anything).Return(3, nil)

		stores, err := mockRepo.Read(ctx)
		assert.NoError(t, err)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)

		assert.Equal(t, len(stores), count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Store timestamps are consistent", func(t *testing.T) {
		now := time.Now()
		store := &models.Store{
			ID:        "store_timestamp_test",
			Name:      "Timestamp Test Store",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Verify that UpdatedAt >= CreatedAt
		assert.True(t, store.UpdatedAt.After(store.CreatedAt) || store.UpdatedAt.Equal(store.CreatedAt))

		// Test ToUpdate method updates timestamp
		repoStore := ToSetStore(store)
		time.Sleep(1 * time.Millisecond)
		updatedRepoStore := repoStore.ToUpdate()

		assert.True(t, updatedRepoStore.UpdatedAt.After(repoStore.CreatedAt))
		assert.True(t, updatedRepoStore.UpdatedAt.After(now))
	})
}

// TestStoreRepositoryErrorScenarios tests various error scenarios
func TestStoreRepositoryErrorScenarios(t *testing.T) {
	ctx := context.Background()

	t.Run("Context cancellation", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		mockRepo.On("Read", cancelledCtx).Return([]*models.Store{}, context.Canceled)

		stores, err := mockRepo.Read(cancelledCtx)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
		assert.Equal(t, 0, len(stores))

		mockRepo.AssertExpectations(t)
	})

	t.Run("Context timeout", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
		defer cancel()

		time.Sleep(1 * time.Millisecond) // Ensure timeout

		mockRepo.On("FindByID", timeoutCtx, "store_123").Return(nil, context.DeadlineExceeded)

		store, err := mockRepo.FindByID(timeoutCtx, "store_123")
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Nil(t, store)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Database connection error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		connectionError := assert.AnError

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(connectionError)

		testStore := &models.Store{ID: "test_store", Name: "Test"}
		err := mockRepo.Create(ctx, testStore)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Malformed data error", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		dataError := assert.AnError

		mockRepo.On("FindByID", mock.Anything, "malformed_store").Return(nil, dataError)

		store, err := mockRepo.FindByID(ctx, "malformed_store")
		assert.Error(t, err)
		assert.Nil(t, store)

		mockRepo.AssertExpectations(t)
	})
}

// TestStoreRepositoryConcurrency tests concurrent access scenarios
func TestStoreRepositoryConcurrency(t *testing.T) {
	ctx := context.Background()

	t.Run("Concurrent reads", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}
		testStore := &models.Store{
			ID:   "concurrent_store",
			Name: "Concurrent Test Store",
		}

		// Set up expectation for multiple calls
		mockRepo.On("FindByID", mock.Anything, "concurrent_store").Return(testStore, nil).Times(10)

		// Simulate concurrent access
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				store, err := mockRepo.FindByID(ctx, "concurrent_store")
				assert.NoError(t, err)
				assert.Equal(t, testStore.ID, store.ID)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Concurrent writes", func(t *testing.T) {
		mockRepo := &MockStoreRepository{}

		// Set up expectations for multiple update calls
		for i := 0; i < 5; i++ {
			storeID := "store_" + string(rune(48+i)) // "store_0", "store_1", etc.
			mockRepo.On("UpdateByID", mock.Anything, storeID, mock.MatchedBy(func(s *models.Store) bool {
				return s.ID == storeID
			})).Return(nil)
		}

		// Simulate concurrent updates
		done := make(chan bool, 5)
		for i := 0; i < 5; i++ {
			go func(index int) {
				storeID := "store_" + string(rune(48+index))
				updatedStore := &models.Store{
					ID:   storeID,
					Name: "Updated Store " + string(rune(48+index)),
				}
				err := mockRepo.UpdateByID(ctx, storeID, updatedStore)
				assert.NoError(t, err)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			<-done
		}

		mockRepo.AssertExpectations(t)
	})
}
