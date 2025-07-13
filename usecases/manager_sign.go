package usecases

import (
	"backend/models"
	"backend/repositories"
	"context"
)

func (u *UseCase) ManagerSignUp(ctx context.Context, email, password string) error {
	// 入力値からマネージャーを作成
	manager := models.NewManager(email, password)

	// パスワードのハッシュ化
	if err := manager.ToEncryptPassword(); err != nil {
		return err
	}

	return repositories.NewManagerRepository(u.db).Create(ctx, manager)
}

func (u *UseCase) ManagerSignIn(ctx context.Context, email, password string) error {
	gotUser, err := repositories.NewManagerRepository(u.db).FindByID(ctx, email)
	if err != nil {
		return err
	}

	// パスワードの検証
	return gotUser.IsVerifyPassword(password)
}
