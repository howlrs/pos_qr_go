package usecases

import (
	"cloud.google.com/go/firestore"
)

type UseCase struct {
	db *firestore.Client
}

func New(db *firestore.Client) *UseCase {
	return &UseCase{
		db: db,
	}
}
