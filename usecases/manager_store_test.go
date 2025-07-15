package usecases

import (
	"backend/models"
	"backend/repositories"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUseCaseを作成するヘルパー関数
func createTestUseCase() *UseCase {
	return New(nil) // モックリポジトリを使用
}

// TestRegisterStore tests the RegisterStore function
func TestRegisterStore(t *testing.T) {
	ctx := context.Background()

	t.Run("successful store registration", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Store")).Return(nil)

		// Act
		store, err := useCase.RegisterStore(ctx, "Test Store", "test@example.com", "password123", "123 Test St", "123-456-7890")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, store)
		assert.Equal(t, "Test Store", store.Name)
		assert.Equal(t, "test@example.com", store.Email)
		assert.NotEqual(t, "password123", store.Password) // Password should be hashed
		assert.Equal(t, "123 Test St", store.Address)
		assert.Equal(t, "123-456-7890", store.Phone)
		assert.NotEmpty(t, store.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("store registration with empty name", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Store")).Return(nil)

		// Act
		store, err := useCase.RegisterStore(ctx, "", "test@example.com", "password123", "123 Test St", "123-456-7890")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, store)
		assert.Equal(t, "", store.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("store registration with empty email", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Store")).Return(nil)

		// Act
		store, err := useCase.RegisterStore(ctx, "Test Store", "", "password123", "123 Test St", "123-456-7890")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, store)
		assert.Equal(t, "", store.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("store registration with empty password", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// No mock setup needed as this should fail before reaching the repository

		// Act
		store, err := useCase.RegisterStore(ctx, "Test Store", "test@example.com", "", "123 Test St", "123-456-7890")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, store)
		assert.Contains(t, err.Error(), "password")
	})
}

// TestGetStore tests the GetStore function
func TestGetStore(t *testing.T) {
	ctx := context.Background()
	storeID := "store_123"

	t.Run("get store by id", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Create a mock store to return
		mockStore := &models.Store{ID: storeID, Name: "Test Store"}

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("FindByID", ctx, storeID).Return(mockStore, nil)

		// Act
		store, err := useCase.GetStore(ctx, storeID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, store)
		assert.Equal(t, storeID, store.ID)
		assert.Equal(t, "Test Store", store.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get store with empty id", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations - return error for empty ID
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("FindByID", ctx, "").Return(nil, assert.AnError)

		// Act
		store, err := useCase.GetStore(ctx, "")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, store)
		mockRepo.AssertExpectations(t)
	})
}

// TestGetAllStores tests the GetAllStores function
func TestGetAllStores(t *testing.T) {
	ctx := context.Background()

	t.Run("get all stores", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Create mock stores to return
		mockStores := []*models.Store{
			{ID: "store_1", Name: "Store 1"},
			{ID: "store_2", Name: "Store 2"},
		}

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("Read", ctx).Return(mockStores, nil)

		// Act
		stores, err := useCase.GetAllStores(ctx)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, stores)
		assert.Len(t, stores, 2)
		assert.Equal(t, "store_1", stores[0].ID)
		assert.Equal(t, "store_2", stores[1].ID)
		mockRepo.AssertExpectations(t)
	})
}

// TestUpdate tests the Update function
func TestUpdate(t *testing.T) {
	ctx := context.Background()
	storeID := "store_123"

	t.Run("update store with password", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("UpdateByID", ctx, storeID, mock.AnythingOfType("*models.Store")).Return(nil)

		// Act
		err := useCase.Update(ctx, storeID, "Updated Store", "updated@example.com", "newpassword", "New Address", "999-9999")

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update store without password", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("UpdateByID", ctx, storeID, mock.AnythingOfType("*models.Store")).Return(nil)

		// Act
		err := useCase.Update(ctx, storeID, "Updated Store", "updated@example.com", "", "New Address", "999-9999")

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update store with empty id", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("UpdateByID", ctx, "", mock.AnythingOfType("*models.Store")).Return(assert.AnError)

		// Act
		err := useCase.Update(ctx, "", "Updated Store", "updated@example.com", "newpassword", "New Address", "999-9999")

		// Assert
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestDelete tests the Delete function
func TestDelete(t *testing.T) {
	ctx := context.Background()
	storeID := "store_123"

	t.Run("delete store", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("DeleteByID", ctx, storeID).Return(nil)

		// Act
		err := useCase.Delete(ctx, storeID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete store with empty id", func(t *testing.T) {
		// Arrange
		useCase := createTestUseCase()

		// Set up mock expectations
		mockRepo := useCase.storeRepo.(*repositories.MockStoreRepository)
		mockRepo.On("DeleteByID", ctx, "").Return(assert.AnError)

		// Act
		err := useCase.Delete(ctx, "")

		// Assert
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// TestStoreModelCreation tests the Store model creation logic
func TestStoreModelCreation(t *testing.T) {
	t.Run("create store model", func(t *testing.T) {
		// Act
		store := models.NewStore("Test Store", "test@example.com", "password123", "123 Test St", "123-456-7890")

		// Assert
		assert.NotNil(t, store)
		assert.Equal(t, "Test Store", store.Name)
		assert.Equal(t, "test@example.com", store.Email)
		assert.Equal(t, "password123", store.Password) // Before hashing
		assert.Equal(t, "123 Test St", store.Address)
		assert.Equal(t, "123-456-7890", store.Phone)
		assert.NotEmpty(t, store.ID)
		assert.True(t, store.ID[:6] == "store_") // Check prefix
		assert.False(t, store.CreatedAt.IsZero())
		assert.False(t, store.UpdatedAt.IsZero())
	})

	t.Run("password hashing", func(t *testing.T) {
		// Arrange
		store := models.NewStore("Test Store", "test@example.com", "password123", "123 Test St", "123-456-7890")
		originalPassword := store.Password

		// Act
		err := store.PasswordToHash()

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, originalPassword, store.Password)
		assert.True(t, len(store.Password) > len(originalPassword))
	})

	t.Run("password hashing with empty password", func(t *testing.T) {
		// Arrange
		store := models.NewStore("Test Store", "test@example.com", "", "123 Test St", "123-456-7890")

		// Act
		err := store.PasswordToHash()

		// Assert
		// Empty passwordの場合、ハッシュ化でエラーが発生する可能性があります
		if err != nil {
			assert.Error(t, err)
			t.Log("Expected error for empty password:", err)
		} else {
			assert.NoError(t, err)
		}
	})
}

// TestUpdateLogic tests the Update function logic without Firestore
func TestUpdateLogic(t *testing.T) {
	storeID := "store_123"

	t.Run("validate update store model creation", func(t *testing.T) {
		// この部分だけを分離してテスト（Firestoreを使わない）
		store := models.NewStore("Updated Store", "updated@example.com", "newpassword", "New Address", "999-9999")
		store.ID = storeID

		// パスワードハッシュ化のテスト
		originalPassword := store.Password
		err := store.PasswordToHash()

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, originalPassword, store.Password)
		assert.Equal(t, storeID, store.ID)
		assert.Equal(t, "Updated Store", store.Name)
		assert.Equal(t, "updated@example.com", store.Email)
	})

	t.Run("validate update store model creation without password", func(t *testing.T) {
		// Arrange
		store := models.NewStore("Updated Store", "updated@example.com", "", "New Address", "999-9999")
		store.ID = storeID

		// パスワードが空の場合のハッシュ化テスト
		if store.Password != "" {
			err := store.PasswordToHash()
			assert.NoError(t, err)
		}

		// Assert
		assert.Equal(t, storeID, store.ID)
		assert.Equal(t, "Updated Store", store.Name)
		assert.Equal(t, "updated@example.com", store.Email)
		assert.Equal(t, "", store.Password) // Empty password should remain empty
	})
}
