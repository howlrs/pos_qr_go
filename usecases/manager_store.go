package usecases

import (
	"backend/models"
	"context"
	"fmt"
)

func (u *UseCase) RegisterStore(ctx context.Context, name, email, password, address, phone string) (*models.Store, error) {
	// 入力値から店舗を作成
	store := models.NewStore(name, email, password, address, phone)

	if err := store.PasswordToHash(); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// レポジトリ層を使用
	if err := u.storeRepo.Create(ctx, store); err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	return store, nil
}

func (u *UseCase) GetStore(ctx context.Context, id string) (*models.Store, error) {
	// レポジトリ層を使用
	return u.storeRepo.FindByID(ctx, id)
}

func (u *UseCase) GetAllStores(ctx context.Context) ([]*models.Store, error) {
	// レポジトリ層を使用
	stores, err := u.storeRepo.Read(ctx)
	if err != nil {
		return nil, err
	}
	if len(stores) == 0 {
		return nil, fmt.Errorf("no stores found")
	}

	return stores, nil
}

func (u *UseCase) Update(ctx context.Context, id, name, email, password, address, phone string) error {
	store := models.NewStore(name, email, password, address, phone)
	store.ID = id

	// パスワードがあれば、ハッシュ化する
	// なければ登録されない
	if store.Password != "" {
		if err := store.PasswordToHash(); err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
	}

	return u.storeRepo.UpdateByID(ctx, store.ID, store)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.storeRepo.DeleteByID(ctx, id)
}
