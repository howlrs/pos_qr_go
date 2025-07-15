package repositories

import (
	"backend/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewSessionRepository tests the NewSessionRepository function
func TestNewSessionRepository(t *testing.T) {
	t.Run("NewSessionRepository with nil client returns MockSessionRepository", func(t *testing.T) {
		repo := NewSessionRepository(nil)
		assert.NotNil(t, repo)

		// Check if it's a mock repository
		_, ok := repo.(*MockSessionRepository)
		assert.True(t, ok, "Should return a MockSessionRepository when client is nil")
	})
}

// TestMockSessionRepository tests the MockSessionRepository implementation
func TestMockSessionRepository(t *testing.T) {
	ctx := context.Background()

	// Test data
	now := time.Now()
	testOrder := models.Order{
		OrderID:   "order_123",
		ProductID: "product_456",
		Quantity:  2,
		Price:     1000.0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	testSession := &models.Session{
		ID:          "session_test123",
		StoreID:     "store_456",
		SeatID:      "seat_789",
		Items:       []models.Order{testOrder},
		TotalAmount: 2000.0,
		Status:      models.StatusCreated,
		ExpiresAt:   now.Add(15 * time.Minute),
		IssuedAt:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	testSessions := []*models.Session{testSession}

	t.Run("Create", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Create", mock.Anything, testSession).Return(nil)

		err := mockRepo.Create(ctx, testSession)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Read", mock.Anything).Return(testSessions, nil)

		sessions, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(sessions))
		assert.Len(t, sessions, 1)
		assert.Equal(t, testSession.ID, sessions[0].ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("FindByID", mock.Anything, "session_test123").Return(testSession, nil)

		session, err := mockRepo.FindByID(ctx, "session_test123")
		assert.NoError(t, err)
		assert.Equal(t, testSession.ID, session.ID)
		assert.Equal(t, testSession.StoreID, session.StoreID)
		assert.Equal(t, testSession.SeatID, session.SeatID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("FindByID", mock.Anything, "nonexistent").Return(nil, assert.AnError)

		session, err := mockRepo.FindByID(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, session)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("FindByField", mock.Anything, "store_id", "store_456").Return(testSessions, nil)

		sessions, err := mockRepo.FindByField(ctx, "store_id", "store_456")
		assert.NoError(t, err)
		assert.Len(t, sessions, 1)
		assert.Equal(t, testSession.StoreID, sessions[0].StoreID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField not found", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("FindByField", mock.Anything, "store_id", "nonexistent_store").Return([]*models.Session{}, nil)

		sessions, err := mockRepo.FindByField(ctx, "store_id", "nonexistent_store")
		assert.NoError(t, err)
		assert.Len(t, sessions, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		updatedSession := &models.Session{
			ID:          "session_test123",
			StoreID:     "store_456",
			SeatID:      "seat_789",
			Items:       testSession.Items,
			TotalAmount: 2500.0,
			Status:      models.StatusConfirmed,
			ExpiresAt:   testSession.ExpiresAt,
			IssuedAt:    testSession.IssuedAt,
			CreatedAt:   testSession.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("UpdateByID", mock.Anything, "session_test123", updatedSession).Return(nil)

		err := mockRepo.UpdateByID(ctx, "session_test123", updatedSession)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "session_test123").Return(nil)

		err := mockRepo.DeleteByID(ctx, "session_test123")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Count", mock.Anything).Return(5, nil)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count error", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Count", mock.Anything).Return(0, assert.AnError)

		count, err := mockRepo.Count(ctx)
		assert.Error(t, err)
		assert.Equal(t, 0, count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists true", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Exists", mock.Anything, "session_test123").Return(true, nil)

		exists, err := mockRepo.Exists(ctx, "session_test123")
		assert.NoError(t, err)
		assert.True(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists false", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Exists", mock.Anything, "nonexistent").Return(false, nil)

		exists, err := mockRepo.Exists(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Exists error", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Exists", mock.Anything, "error_id").Return(false, assert.AnError)

		exists, err := mockRepo.Exists(ctx, "error_id")
		assert.Error(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})
}

// TestSessionRepositoryInterface tests that SessionRepository implements the Repository interface
func TestSessionRepositoryInterface(t *testing.T) {
	t.Run("SessionRepository implements Repository interface", func(t *testing.T) {
		var repo Repository[models.Session]
		repo = NewSessionRepository(nil) // This should return a MockSessionRepository
		assert.NotNil(t, repo)

		// Verify the repository is a MockSessionRepository
		_, ok := repo.(*MockSessionRepository)
		assert.True(t, ok, "Should return a MockSessionRepository when client is nil")
	})
}

// TestRepositoryErrorHandling tests error handling scenarios
func TestSessionRepositoryErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("Create error", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		testSession := &models.Session{ID: "test_id"}
		mockRepo.On("Create", mock.Anything, testSession).Return(assert.AnError)

		err := mockRepo.Create(ctx, testSession)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Read error", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Read", mock.Anything).Return([]*models.Session{}, assert.AnError)

		sessions, err := mockRepo.Read(ctx)
		assert.Error(t, err)
		assert.Len(t, sessions, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateByID error", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		testSession := &models.Session{ID: "test_id"}
		mockRepo.On("UpdateByID", mock.Anything, "test_id", testSession).Return(assert.AnError)

		err := mockRepo.UpdateByID(ctx, "test_id", testSession)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteByID error", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("DeleteByID", mock.Anything, "test_id").Return(assert.AnError)

		err := mockRepo.DeleteByID(ctx, "test_id")
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

// TestSessionStruct tests the Session struct and its conversion methods
func TestSessionStruct(t *testing.T) {
	now := time.Now()
	testOrder := models.Order{
		OrderID:   "order_123",
		ProductID: "product_456",
		Quantity:  2,
		Price:     1000.0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	testSession := &models.Session{
		ID:          "session_123",
		StoreID:     "store_456",
		SeatID:      "seat_789",
		Items:       []models.Order{testOrder},
		TotalAmount: 2000.0,
		Status:      models.StatusCreated,
		ExpiresAt:   now.Add(15 * time.Minute),
		IssuedAt:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("ToSetSession conversion", func(t *testing.T) {
		repoSession := ToSetSession(testSession)
		assert.NotNil(t, repoSession)
		assert.Equal(t, testSession.ID, repoSession.ID)
		assert.Equal(t, testSession.StoreID, repoSession.StoreID)
		assert.Equal(t, testSession.SeatID, repoSession.SeatID)
		assert.Equal(t, testSession.TotalAmount, repoSession.TotalAmount)
		assert.Equal(t, Status(testSession.Status), repoSession.Status)
		assert.Equal(t, testSession.ExpiresAt, repoSession.ExpiresAt)
		assert.Equal(t, testSession.IssuedAt, repoSession.IssuedAt)
		assert.Equal(t, testSession.CreatedAt, repoSession.CreatedAt)
		assert.Equal(t, testSession.UpdatedAt, repoSession.UpdatedAt)
		assert.Len(t, repoSession.Items, 1)
	})

	t.Run("ToUpdate method", func(t *testing.T) {
		repoSession := ToSetSession(testSession)
		originalUpdatedAt := repoSession.UpdatedAt

		// Sleep for a small duration to ensure time difference
		time.Sleep(1 * time.Millisecond)

		updatedSession := repoSession.ToUpdate()
		assert.NotNil(t, updatedSession)
		assert.True(t, updatedSession.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("ToModel conversion", func(t *testing.T) {
		repoSession := &Session{
			ID:          "session_456",
			StoreID:     "store_789",
			SeatID:      "seat_123",
			Items:       []Order{{OrderID: "order_456", ProductID: "product_789", Quantity: 1, Price: 500.0, CreatedAt: now, UpdatedAt: now}},
			TotalAmount: 500.0,
			Status:      Status(models.StatusConfirmed),
			ExpiresAt:   now.Add(10 * time.Minute),
			IssuedAt:    now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		modelSession := repoSession.ToModel()
		assert.NotNil(t, modelSession)
		assert.Equal(t, repoSession.ID, modelSession.ID)
		assert.Equal(t, repoSession.StoreID, modelSession.StoreID)
		assert.Equal(t, repoSession.SeatID, modelSession.SeatID)
		assert.Equal(t, repoSession.TotalAmount, modelSession.TotalAmount)
		assert.Equal(t, models.Status(repoSession.Status), modelSession.Status)
		assert.Equal(t, repoSession.ExpiresAt, modelSession.ExpiresAt)
		assert.Equal(t, repoSession.IssuedAt, modelSession.IssuedAt)
		assert.Equal(t, repoSession.CreatedAt, modelSession.CreatedAt)
		assert.Equal(t, repoSession.UpdatedAt, modelSession.UpdatedAt)
		assert.Len(t, modelSession.Items, 1)
	})
}

// TestOrderConversions tests Order conversion methods
func TestOrderConversions(t *testing.T) {
	now := time.Now()
	testOrders := []models.Order{
		{
			OrderID:   "order_1",
			ProductID: "product_1",
			Quantity:  2,
			Price:     1000.0,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			OrderID:   "order_2",
			ProductID: "product_2",
			Quantity:  1,
			Price:     500.0,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	t.Run("ToSetOrders conversion", func(t *testing.T) {
		repoOrders := ToSetOrders(testOrders)
		assert.Len(t, repoOrders, 2)

		for i, repoOrder := range repoOrders {
			assert.Equal(t, testOrders[i].OrderID, repoOrder.OrderID)
			assert.Equal(t, testOrders[i].ProductID, repoOrder.ProductID)
			assert.Equal(t, testOrders[i].Quantity, repoOrder.Quantity)
			assert.Equal(t, testOrders[i].Price, repoOrder.Price)
			assert.Equal(t, testOrders[i].CreatedAt, repoOrder.CreatedAt)
			assert.Equal(t, testOrders[i].UpdatedAt, repoOrder.UpdatedAt)
		}
	})

	t.Run("ToModelOrders conversion", func(t *testing.T) {
		repoOrders := []Order{
			{
				OrderID:   "order_3",
				ProductID: "product_3",
				Quantity:  3,
				Price:     750.0,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		modelOrders := ToModelOrders(repoOrders)
		assert.Len(t, modelOrders, 1)
		assert.Equal(t, repoOrders[0].OrderID, modelOrders[0].OrderID)
		assert.Equal(t, repoOrders[0].ProductID, modelOrders[0].ProductID)
		assert.Equal(t, repoOrders[0].Quantity, modelOrders[0].Quantity)
		assert.Equal(t, repoOrders[0].Price, modelOrders[0].Price)
		assert.Equal(t, repoOrders[0].CreatedAt, modelOrders[0].CreatedAt)
		assert.Equal(t, repoOrders[0].UpdatedAt, modelOrders[0].UpdatedAt)
	})
}

// TestSessionRepositoryBusinessLogic tests business logic scenarios
func TestSessionRepositoryBusinessLogic(t *testing.T) {
	ctx := context.Background()

	t.Run("Find sessions by store", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		now := time.Now()

		storeSessions := []*models.Session{
			{ID: "session_1", StoreID: "store_123", SeatID: "seat_1", Status: models.StatusCreated, CreatedAt: now, UpdatedAt: now},
			{ID: "session_2", StoreID: "store_123", SeatID: "seat_2", Status: models.StatusConfirmed, CreatedAt: now, UpdatedAt: now},
		}

		mockRepo.On("FindByField", mock.Anything, "store_id", "store_123").Return(storeSessions, nil)

		sessions, err := mockRepo.FindByField(ctx, "store_id", "store_123")
		assert.NoError(t, err)
		assert.Len(t, sessions, 2)

		for _, session := range sessions {
			assert.Equal(t, "store_123", session.StoreID)
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Find sessions by seat", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		now := time.Now()

		seatSessions := []*models.Session{
			{ID: "session_3", StoreID: "store_456", SeatID: "seat_789", Status: models.StatusPreparing, CreatedAt: now, UpdatedAt: now},
		}

		mockRepo.On("FindByField", mock.Anything, "seat_id", "seat_789").Return(seatSessions, nil)

		sessions, err := mockRepo.FindByField(ctx, "seat_id", "seat_789")
		assert.NoError(t, err)
		assert.Len(t, sessions, 1)
		assert.Equal(t, "seat_789", sessions[0].SeatID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Find sessions by status", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		now := time.Now()

		activeSessions := []*models.Session{
			{ID: "session_4", StoreID: "store_789", SeatID: "seat_456", Status: models.StatusPreparing, CreatedAt: now, UpdatedAt: now},
			{ID: "session_5", StoreID: "store_789", SeatID: "seat_123", Status: models.StatusPreparing, CreatedAt: now, UpdatedAt: now},
		}

		mockRepo.On("FindByField", mock.Anything, "status", string(models.StatusPreparing)).Return(activeSessions, nil)

		sessions, err := mockRepo.FindByField(ctx, "status", string(models.StatusPreparing))
		assert.NoError(t, err)
		assert.Len(t, sessions, 2)

		for _, session := range sessions {
			assert.Equal(t, models.StatusPreparing, session.Status)
		}

		mockRepo.AssertExpectations(t)
	})
}

// TestSessionRepositoryTransactionScenarios tests transaction-like scenarios
func TestSessionRepositoryTransactionScenarios(t *testing.T) {
	ctx := context.Background()

	// setup is a helper function to initialize a new mock repository for each sub-test.
	setup := func() *MockSessionRepository {
		return &MockSessionRepository{}
	}

	t.Run("Create and immediately verify session exists", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		now := time.Now()
		newSession := &models.Session{
			ID:          "session_new_123",
			StoreID:     "store_456",
			SeatID:      "seat_789",
			Items:       []models.Order{{OrderID: "order_123", ProductID: "product_456", Quantity: 1, Price: 1000.0, CreatedAt: now, UpdatedAt: now}},
			TotalAmount: 1000.0,
			Status:      models.StatusCreated,
			ExpiresAt:   now.Add(15 * time.Minute),
			IssuedAt:    now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mockRepo.On("Create", ctx, newSession).Return(nil)
		mockRepo.On("Exists", ctx, newSession.ID).Return(true, nil)
		mockRepo.On("FindByID", ctx, newSession.ID).Return(newSession, nil)

		// Act
		createErr := mockRepo.Create(ctx, newSession)
		exists, existsErr := mockRepo.Exists(ctx, newSession.ID)
		retrievedSession, findErr := mockRepo.FindByID(ctx, newSession.ID)

		// Assert
		assert.NoError(t, createErr, "Create should not return an error")
		assert.NoError(t, existsErr, "Exists should not return an error")
		assert.True(t, exists, "Session should exist after creation")
		assert.NoError(t, findErr, "FindByID should not return an error")
		assert.Equal(t, newSession, retrievedSession, "Retrieved session should match the created session")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Update session status and verify changes", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		now := time.Now()
		originalSession := &models.Session{
			ID:          "session_update_123",
			StoreID:     "store_456",
			SeatID:      "seat_789",
			Items:       []models.Order{{OrderID: "order_123", ProductID: "product_456", Quantity: 1, Price: 1000.0, CreatedAt: now, UpdatedAt: now}},
			TotalAmount: 1000.0,
			Status:      models.StatusCreated,
			ExpiresAt:   now.Add(15 * time.Minute),
			IssuedAt:    now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		updatedSession := &models.Session{
			ID:          originalSession.ID,
			StoreID:     originalSession.StoreID,
			SeatID:      originalSession.SeatID,
			Items:       originalSession.Items,
			TotalAmount: originalSession.TotalAmount,
			Status:      models.StatusConfirmed,
			ExpiresAt:   originalSession.ExpiresAt,
			IssuedAt:    originalSession.IssuedAt,
			CreatedAt:   originalSession.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("UpdateByID", ctx, originalSession.ID, updatedSession).Return(nil)
		mockRepo.On("FindByID", ctx, originalSession.ID).Return(updatedSession, nil)

		// Act
		updateErr := mockRepo.UpdateByID(ctx, originalSession.ID, updatedSession)
		retrievedSession, findErr := mockRepo.FindByID(ctx, originalSession.ID)

		// Assert
		assert.NoError(t, updateErr, "UpdateByID should not return an error")
		assert.NoError(t, findErr, "FindByID should not return an error after update")
		assert.Equal(t, updatedSession, retrievedSession, "Retrieved session should reflect the updates")
		assert.Equal(t, models.StatusConfirmed, retrievedSession.Status, "Status should be updated")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete session and verify it no longer exists", func(t *testing.T) {
		// Arrange
		mockRepo := setup()
		sessionID := "session_delete_123"

		mockRepo.On("DeleteByID", ctx, sessionID).Return(nil)
		mockRepo.On("Exists", ctx, sessionID).Return(false, nil)
		mockRepo.On("FindByID", ctx, sessionID).Return(nil, assert.AnError)

		// Act
		deleteErr := mockRepo.DeleteByID(ctx, sessionID)
		exists, existsErr := mockRepo.Exists(ctx, sessionID)
		deletedSession, findErr := mockRepo.FindByID(ctx, sessionID)

		// Assert
		assert.NoError(t, deleteErr, "DeleteByID should not return an error")
		assert.NoError(t, existsErr, "Exists should not return an error after deletion")
		assert.False(t, exists, "Session should not exist after deletion")
		assert.Error(t, findErr, "FindByID should return an error for a deleted session")
		assert.Nil(t, deletedSession, "FindByID should return a nil session after deletion")

		mockRepo.AssertExpectations(t)
	})
}

// TestSessionRepositoryEdgeCases tests edge cases and boundary conditions
func TestSessionRepositoryEdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("Read returns empty slice when no sessions exist", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("Read", mock.Anything).Return([]*models.Session{}, nil)

		sessions, err := mockRepo.Read(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, sessions)
		assert.Len(t, sessions, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByField with nil value", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		mockRepo.On("FindByField", mock.Anything, "store_id", nil).Return([]*models.Session{}, nil)

		sessions, err := mockRepo.FindByField(ctx, "store_id", nil)
		assert.NoError(t, err)
		assert.NotNil(t, sessions)
		assert.Len(t, sessions, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Count matches actual number of sessions", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		now := time.Now()
		testSessions := []*models.Session{
			{ID: "session_1", StoreID: "store_123", Status: models.StatusCreated, CreatedAt: now, UpdatedAt: now},
			{ID: "session_2", StoreID: "store_456", Status: models.StatusConfirmed, CreatedAt: now, UpdatedAt: now},
			{ID: "session_3", StoreID: "store_789", Status: models.StatusPreparing, CreatedAt: now, UpdatedAt: now},
		}

		mockRepo.On("Read", mock.Anything).Return(testSessions, nil)
		mockRepo.On("Count", mock.Anything).Return(3, nil)

		sessions, err := mockRepo.Read(ctx)
		assert.NoError(t, err)

		count, err := mockRepo.Count(ctx)
		assert.NoError(t, err)

		assert.Equal(t, len(sessions), count)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Session with expired time", func(t *testing.T) {
		mockRepo := &MockSessionRepository{}
		now := time.Now()
		expiredSession := &models.Session{
			ID:          "session_expired",
			StoreID:     "store_123",
			SeatID:      "seat_456",
			Items:       []models.Order{},
			TotalAmount: 0.0,
			Status:      models.StatusCreated,
			ExpiresAt:   now.Add(-1 * time.Hour), // Already expired
			IssuedAt:    now.Add(-2 * time.Hour),
			CreatedAt:   now.Add(-2 * time.Hour),
			UpdatedAt:   now.Add(-1 * time.Hour),
		}

		mockRepo.On("FindByID", mock.Anything, "session_expired").Return(expiredSession, nil)

		session, err := mockRepo.FindByID(ctx, "session_expired")
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.True(t, session.ExpiresAt.Before(now), "Session should be expired")

		mockRepo.AssertExpectations(t)
	})
}
