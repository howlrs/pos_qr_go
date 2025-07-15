package usecases

import (
	"backend/models"
	"context"
)

func (u *UseCase) ManagerSignUp(ctx context.Context, email, password string) error {
	// 入力値からマネージャーを作成
	manager := models.NewManager(email, password)

	// パスワードのハッシュ化
	if err := manager.ToEncryptPassword(); err != nil {
		return err
	}

	return u.managerRepo.Create(ctx, manager)
}

func (u *UseCase) ManagerSignIn(ctx context.Context, email, password string) error {
	gotUser, err := u.managerRepo.FindByID(ctx, email)
	if err != nil {
		return err
	}

	// パスワードの検証
	return gotUser.IsVerifyPassword(password)
}
