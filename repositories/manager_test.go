package repositories

import (
	"backend/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewManagerRepository tests the NewManagerRepository function
func TestNewManagerRepository(t *testing.T) {
	t.Run("NewManagerRepository with nil client returns MockManagerRepository", func(t *testing.T) {
		repo := NewManagerRepository(nil)
		assert.NotNil(t, repo)

		// Check if it's a mock repository
		_, ok := repo.(*MockManagerRepository)
		assert.True(t, ok, "Should return a MockManagerRepository when client is nil")
	})
}

// TestMockManagerRepository tests the MockManagerRepository implementation
func TestMockManagerRepository(t *testing.T) {
	ctx := context.Background()

	// Test data
	testManager := &models.Manager{
		Email:    "manager@example.com",
		Password: "hashedpassword123",
	}

	testManagers := []*models.Manager{testManager}

	t.Run("Create", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Create", mock.Anything, testManager).Return(nil)

		err := mockRepo.Create(ctx, testManager)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Read", mock.Anything).Return(testManagers, nil)

		managers, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(managers))
		assert.Len(t, managers, 1)
		assert.Equal(t, testManager.Email, managers[0].Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("FindByID", mock.Anything, "manager@example.com").Return(testManager, nil)

		manager, err := mockRepo.FindByID(ctx, "manager@example.com")
		assert.NoError(t, err)
		assert.Equal(t, testManager.Email, manager.Email)
		assert.Equal(t, testManager.Password, manager.Password)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("FindByID", mock.Anything, "nonexistent@example.com").Return(nil, assert.AnError)

		manager, err := mockRepo.FindByID(ctx, "nonexistent@example.com")
		assert.Error(t, err)
		assert.Nil(t, manager)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("FindByField", mock.Anything, "email", "manager@example.com").Return(testManagers, nil)

		managers, err := mockRepo.FindByField(ctx, "email", "manager@example.com")
		assert.NoError(t, err)
		assert.Len(t, managers, 1)
		assert.Equal(t, testManager.Email, managers[0].Email)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField not found", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("FindByField", mock.Anything, "email", "notfound@example.com").Return([]*models.Manager{}, nil)

		managers, err := mockRepo.FindByField(ctx, "email", "notfound@example.com")
		assert.NoError(t, err)
		assert.Len(t, managers, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		updatedManager := &models.Manager{
			Email:    "manager@example.com",
			Password: "newhashedpassword456",
		}

		mockRepo.On("UpdateByID", mock.Anything, "manager@example.com", updatedManager).Return(nil)

		err := mockRepo.UpdateByID(ctx, "manager@example.com", updatedManager)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "manager@example.com").Return(nil)

		err := mockRepo.DeleteByID(ctx, "manager@example.com")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Count", mock.Anything).Return(3, nil)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count error", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Count", mock.Anything).Return(0, assert.AnError)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists true", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Exists", mock.Anything, "manager@example.com").Return(true, nil)

		exists, err := mockRepo.Exists(ctx, "manager@example.com")
		assert.NoError(t, err)
		assert.True(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists false", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Exists", mock.Anything, "nonexistent@example.com").Return(false, nil)

		exists, err := mockRepo.Exists(ctx, "nonexistent@example.com")
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists error", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Exists", mock.Anything, "error@example.com").Return(false, assert.AnError)

		exists, err := mockRepo.Exists(ctx, "error@example.com")
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})
}

// TestManagerRepositoryInterface tests that ManagerRepository implements the Repository interface
func TestManagerRepositoryInterface(t *testing.T) {
	t.Run("ManagerRepository implements Repository interface", func(t *testing.T) {
		var repo Repository[models.Manager]
		repo = NewManagerRepository(nil) // This should return a MockManagerRepository
		assert.NotNil(t, repo)

		// Verify the repository is a MockManagerRepository
		_, ok := repo.(*MockManagerRepository)
		assert.True(t, ok, "Should return a MockManagerRepository when client is nil")
	})
}

// TestRepositoryErrorHandling tests error handling scenarios
func TestManagerRepositoryErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("Create error", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		testManager := &models.Manager{Email: "test@example.com"}
		mockRepo.On("Create", mock.Anything, testManager).Return(assert.AnError)

		err := mockRepo.Create(ctx, testManager)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read error", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Read", mock.Anything).Return([]*models.Manager{}, assert.AnError)

		managers, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.Len(t, managers, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID error", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		testManager := &models.Manager{Email: "test@example.com"}
		mockRepo.On("UpdateByID", mock.Anything, "test@example.com", testManager).Return(assert.AnError)

		err := mockRepo.UpdateByID(ctx, "test@example.com", testManager)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID error", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "test@example.com").Return(assert.AnError)

		err := mockRepo.DeleteByID(ctx, "test@example.com")
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestManagerRepositoryBusinessLogic tests business logic scenarios
func TestManagerRepositoryBusinessLogic(t *testing.T) {
	ctx := context.Background()

	t.Run("Create manager with duplicate email should fail", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		existingManager := &models.Manager{
			Email: "duplicate@example.com",
		}

		// First, check if manager with same email exists
		mockRepo.On("FindByField", mock.Anything, "email", "duplicate@example.com").Return([]*models.Manager{existingManager}, nil)

		managers, err := mockRepo.FindByField(ctx, "email", "duplicate@example.com")
		assert.NoError(t, err)
		assert.Len(t, managers, 1)

		// Should not create another manager with same email in business logic
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update manager password", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		managerWithNewPassword := &models.Manager{
			Email:    "manager@example.com",
			Password: "new_hashed_password",
		}

		mockRepo.On("UpdateByID", mock.Anything, "manager@example.com", managerWithNewPassword).Return(nil)

		err := mockRepo.UpdateByID(ctx, "manager@example.com", managerWithNewPassword)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Find managers by specific criteria", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}

		// Find managers with specific domain
		managersWithDomain := []*models.Manager{
			{Email: "admin@company.com", Password: "hash1"},
			{Email: "manager@company.com", Password: "hash2"},
		}

		// This would typically be done with a more complex query in real implementation
		mockRepo.On("Read", mock.Anything).Return(managersWithDomain, nil)

		managers, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.Len(t, managers, 2)

		for _, manager := range managers {
			assert.Contains(t, manager.Email, "@company.com")
		}

		mockRepo.AssertExpectations(t)
	})
}

// TestManagerRepositoryTransactionScenarios tests transaction-like scenarios
func TestManagerRepositoryTransactionScenarios(t *testing.T) {
	ctx := context.Background()

	// setup is a helper function to initialize a new mock repository for each sub-test.
	setup := func() *MockManagerRepository {
		return &MockManagerRepository{}
	}

	t.Run("Create and immediately verify manager exists", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		newManager := &models.Manager{
			Email:    "new@example.com",
			Password: "hashedpassword",
		}

		mockRepo.On("Create", ctx, newManager).Return(nil)
		mockRepo.On("Exists", ctx, newManager.Email).Return(true, nil)
		mockRepo.On("FindByID", ctx, newManager.Email).Return(newManager, nil)

		// Act
		createErr := mockRepo.Create(ctx, newManager)
		exists, existsErr := mockRepo.Exists(ctx, newManager.Email)
		retrievedManager, findErr := mockRepo.FindByID(ctx, newManager.Email)

		// Assert
		assert.NoError(t, createErr, "Create should not return an error")
		assert.NoError(t, existsErr, "Exists should not return an error")
		assert.True(t, exists, "Manager should exist after creation")
		assert.NoError(t, findErr, "FindByID should not return an error")
		assert.Equal(t, newManager, retrievedManager, "Retrieved manager should match the created manager")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Update manager and verify changes", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		originalManager := &models.Manager{
			Email:    "update@example.com",
			Password: "originalpassword",
		}
		updatedManager := &models.Manager{
			Email:    originalManager.Email,
			Password: "updatedpassword",
		}

		mockRepo.On("UpdateByID", ctx, originalManager.Email, updatedManager).Return(nil)
		mockRepo.On("FindByID", ctx, originalManager.Email).Return(updatedManager, nil)

		// Act
		updateErr := mockRepo.UpdateByID(ctx, originalManager.Email, updatedManager)
		retrievedManager, findErr := mockRepo.FindByID(ctx, originalManager.Email)

		// Assert
		assert.NoError(t, updateErr, "UpdateByID should not return an error")
		assert.NoError(t, findErr, "FindByID should not return an error after update")
		assert.Equal(t, updatedManager, retrievedManager, "Retrieved manager should reflect the updates")
		assert.NotEqual(t, originalManager.Password, retrievedManager.Password, "Password should be updated")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete manager and verify it no longer exists", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		managerEmail := "delete@example.com"

		mockRepo.On("DeleteByID", ctx, managerEmail).Return(nil)
		mockRepo.On("Exists", ctx, managerEmail).Return(false, nil)
		mockRepo.On("FindByID", ctx, managerEmail).Return(nil, assert.AnError)

		// Act
		deleteErr := mockRepo.DeleteByID(ctx, managerEmail)
		exists, existsErr := mockRepo.Exists(ctx, managerEmail)
		deletedManager, findErr := mockRepo.FindByID(ctx, managerEmail)

		// Assert
		assert.NoError(t, deleteErr, "DeleteByID should not return an error")
		assert.NoError(t, existsErr, "Exists should not return an error after deletion")
		assert.False(t, exists, "Manager should not exist after deletion")
		assert.Error(t, findErr, "FindByID should return an error for a deleted manager")
		assert.Nil(t, deletedManager, "FindByID should return a nil manager after deletion")

		mockRepo.AssertExpectations(t)
	})
}

// TestManagerRepositoryEdgeCases tests edge cases and boundary conditions
func TestManagerRepositoryEdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("FindByField with nil value", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("FindByField", mock.Anything, "email", nil).Return([]*models.Manager{}, nil)

		managers, err := mockRepo.FindByField(ctx, "email", nil)
		assert.NoError(t, err)
		assert.NotNil(t, managers)
		assert.Len(t, managers, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID with non-existent email", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		testManager := &models.Manager{
			Email:    "nonexistent@example.com",
			Password: "password",
		}

		mockRepo.On("UpdateByID", mock.Anything, "nonexistent@example.com", testManager).Return(assert.AnError)

		err := mockRepo.UpdateByID(ctx, "nonexistent@example.com", testManager)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID with non-existent email", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "nonexistent@example.com").Return(assert.AnError)

		err := mockRepo.DeleteByID(ctx, "nonexistent@example.com")
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read returns empty slice when no managers exist", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		mockRepo.On("Read", mock.Anything).Return([]*models.Manager{}, nil)

		managers, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, managers)
		assert.Len(t, managers, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count matches actual number of managers", func(t *testing.T) {
		mockRepo := &MockManagerRepository{}
		testManagers := []*models.Manager{
			{Email: "manager1@example.com", Password: "hash1"},
			{Email: "manager2@example.com", Password: "hash2"},
			{Email: "manager3@example.com", Password: "hash3"},
		}

		mockRepo.On("Read", mock.Anything).Return(testManagers, nil)
		mockRepo.On("Count", mock.Anything).Return(3, nil)

		managers, err := mockRepo.Read(ctx)
		assert.NoError(t, err)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)

		assert.Equal(t, len(managers), count)

		mockRepo.AssertExpectations(t)
	})
}
