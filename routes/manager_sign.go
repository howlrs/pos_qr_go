package routes

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	manager := &models.Manager{}
	if err := c.Bind(manager); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
	}

	// パスワードのハッシュ化
	if err := manager.ToEncryptPassword(); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to encrypt password")
	}

	// set database
	manager.ID = xid.New().String()
	if _, err := p.firestore.Collection(manager.ToCollection(p.IsTest())).Doc(manager.ID).Set(c.Request().Context(), manager); err != nil {
		// すでにKeyが存在する場合はエラーを返す
		if status.Code(err) == codes.AlreadyExists {
			return responseHandler(c, http.StatusConflict, nil, err, "Failed to set manager, already exists")
		}

		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to set manager")
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
	manager := &models.Manager{}
	dbManager := &models.Manager{}
	if !p.IsTest() {
		// リクエストのバインド
		if err := c.Bind(manager); err != nil {
			return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind request")
		}
		if manager.Email == "" || manager.Password == "" {
			return responseHandler(c, http.StatusBadRequest, nil, nil, "Failed to bind request")
		}

		// read database
		// KeyはEmailを想定
		doc, err := p.firestore.Collection(manager.ToCollection(p.IsTest())).Doc(manager.Email).Get(c.Request().Context())
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return responseHandler(c, http.StatusNotFound, nil, err, "Failed to get manager, not exists")
			}

			return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get manager")
		}

		if err := doc.DataTo(dbManager); err != nil {
			return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get manager")
		}

		// パスワードの検証
		if err := dbManager.IsVerifyPassword(manager.Password); err != nil {
			return responseHandler(c, http.StatusUnauthorized, nil, err, "Failed to verify password")
		}

		// [Important] パスワードは返さない
		dbManager.Password = ""
	} else {
		// テスト用のユーザ情報を設定
		dbManager = &models.Manager{
			ID:       "test",
			Email:    "",
			Password: "test",
		}
	}

	// 期限を指定し
	// ユーザ情報からJWTトークンを生成
	expireAt := time.Now().Add(time.Hour * 24 * 7)
	isAdmin := false
	token, err := models.NewClaims(dbManager, isAdmin, expireAt).ToJwtToken()
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
		"manager":    dbManager,
	}, nil, "success, create jwt token")
}
