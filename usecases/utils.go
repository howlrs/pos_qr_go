package usecases

import (
	"backend/models"
	"backend/repositories"

	"cloud.google.com/go/firestore"
)

type UseCase struct {
	managerRepo repositories.Repository[models.Manager]
	sessionRepo repositories.Repository[models.Session]
	seatRepo    repositories.Repository[models.Seat]
	storeRepo   repositories.Repository[models.Store]
}

func New(db *firestore.Client) *UseCase {
	return &UseCase{
		managerRepo: repositories.NewManagerRepository(db),
		sessionRepo: repositories.NewSessionRepository(db),
		seatRepo:    repositories.NewSeatRepository(db),
		storeRepo:   repositories.NewStoreRepository(db),
	}
}
