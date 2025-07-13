package routes

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestManager struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *RequestManager) IsValidate() error {
	if m.Email == "" || m.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email and Password are required")
	}
	return nil
}

// Signup は新しいマネージャーアカウントを作成するためのハンドラです。
// リクエストボディからマネージャー情報をバインドし、パスワードをハッシュ化した上で
// Firestore に保存します。保存時に既存のIDが存在する場合は409 Conflictを返します。
//
// 引数:
//
//	c echo.Context - リクエストコンテキスト
//
// 戻り値:
//
//	error - HTTPレスポンスを返します。成功時は200 OK、バリデーションエラー時は400 Bad Request、
//	       パスワードハッシュ化失敗時やDBエラー時は500 Internal Server Error、
//	       既存IDの場合は409 Conflictを返します。
func (p *Client) Signup(c echo.Context) error {
	// リクエストのバインド
	manager := &RequestManager{}
	if err := c.Bind(manager); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
	}

	// 入力値の検証
	if err := manager.IsValidate(); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to validate request")
	}

	// データベースへ保存
	// - パスワードのハッシュ化
	if err := p.uc.ManagerSignUp(c.Request().Context(), manager.Email, manager.Email); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to sign up manager")
	}

	return responseHandler(c, http.StatusOK, manager, nil, "success, create manager")
}

// Signin handles manager sign-in requests.
//
// 引数:
//   - c (echo.Context): Echoのコンテキスト。リクエスト情報やレスポンスの送信に使用されます。
//
// 返り値:
//   - error: エラーが発生した場合は適切なHTTPステータスとエラーメッセージを返します。正常時はJWTトークンとマネージャ情報を含むレスポンスを返します。
//
// 概要:
//   - リクエストボディからマネージャ情報（Email, Password）をバインドします。
//   - Firestoreから該当マネージャ情報を取得し、パスワードを検証します。
//   - パスワードが一致した場合、JWTトークンを生成して返却します。
//   - テストモードの場合はテスト用のマネージャ情報を使用します。
//   - パスワードはレスポンスに含めません。
func (p *Client) Signin(c echo.Context) error {
	manager := &RequestManager{}
	if !p.IsTest() {
		// リクエストのバインド
		if err := c.Bind(manager); err != nil {
			return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
		}
		if err := manager.IsValidate(); err != nil {
			return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to validate request")
		}

		// read database
		// KeyはEmailを想定
		if err := p.uc.ManagerSignIn(c.Request().Context(), manager.Email, manager.Password); err != nil {
			return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to sign in manager")
		}

		// [Important] パスワードは返さない
		manager.Password = ""
	} else {
		// テスト用のユーザ情報を設定
		manager = &RequestManager{
			Email:    "",
			Password: "",
		}
	}

	setManager := &models.Manager{
		Email:    manager.Email,
		Password: manager.Password,
	}

	// 期限を指定し
	// ユーザ情報からJWTトークンを生成
	expireAt := time.Now().Add(time.Hour * 24 * 7)
	isAdmin := false
	token, err := models.NewClaims(setManager, isAdmin, expireAt).ToJwtToken()
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to create token")
	}

	c.SetCookie(&http.Cookie{
		Name:    "jwt_token",
		Value:   token,
		Expires: expireAt,
		// HttpOnly: true, // XSS攻撃から保護
		// Secure:   !p.IsTest(),          // 開発環境ではfalse、本番環境ではtrue
		// SameSite: http.SameSiteLaxMode, // CSRF攻撃から保護
	})

	return responseHandler(c, http.StatusOK, echo.Map{
		"token":      token,
		"token_type": "bearer",
		"manager":    manager,
	}, nil, "success, create jwt token")
}
