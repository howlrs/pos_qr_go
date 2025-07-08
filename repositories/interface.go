package repositories

import "context"

type Repository[T any] interface {
	Create(ctx context.Context, data *T) error

	Read(ctx context.Context) ([]*T, error)
	FindByID(ctx context.Context, id string) (*T, error)
	FindByField(ctx context.Context, field string, value any) ([]*T, error)

	UpdateByID(ctx context.Context, id string, data *T) error
	DeleteByID(ctx context.Context, id string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)
}
