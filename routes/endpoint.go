package routes

import (
	"backend/models"
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Client struct {
	isTest    bool
	firestore *firestore.Client
}

func NewClient(isTest bool) *Client {
	ctx := context.Background()
	projectId := os.Getenv("PROJECT_ID")
	db, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		panic(err)
	}

	return &Client{
		isTest:    isTest,
		firestore: db,
	}
}

// IsTest: テストモードかどうかを返す
// コレクション名を切り替える際やテスト用のデータを返す際に使用
func (p *Client) IsTest() bool {
	return p.isTest
}

func setmiddleware(isTest bool) echo.MiddlewareFunc {
	var setCors = middleware.CORS()
	if !isTest { // 本番環境は展開先のGUIのURLを設定
		origins := []string{
			os.Getenv("FRONTEND_URL"),
		}

		setCors = middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: origins,
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
				echo.HeaderAuthorization,
				echo.HeaderXRequestedWith,
			},
			AllowMethods: []string{
				echo.GET,
				echo.POST,
				echo.PUT,
				echo.DELETE,
				echo.OPTIONS,
			},
			MaxAge: 86400,
		})
	}

	return setCors
}

// Endpoint sets up the routes for the application.
func Endpoint(e *echo.Echo, isTest bool) {
	p := NewClient(isTest)

	e.Use(setmiddleware(isTest))

	// version 1
	v1 := e.Group("/api/v1")

	// public routes
	v1Public := v1.Group("/public")
	v1Public.GET("/health", publicHealth)
	v1Public.POST("/signup", p.Signup)
	v1Public.POST("/signin", p.Signin)

	v1Private := v1.Group("/private")
	jwtSecret := os.Getenv("JWT_SECRET")

	// private routes
	// Manager用のログインJWTトークンを使用するためのルート
	p.handleManager(v1Private, jwtSecret)

	// private session routes
	// 注文セッション用のJWTトークンを使用するためのルート
	p.handleSession(v1Private, jwtSecret)
}

func publicHealth(c echo.Context) error {
	return responseHandler(c, http.StatusOK, echo.Map{"message": "OK"}, nil, "success, public health")
}

func privateHealth(c echo.Context) error {
	// コンテキストからユーザ情報を取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.Claims)
	fmt.Println(claims)

	return responseHandler(c, http.StatusOK, echo.Map{"message": "OK"}, nil, "success, private health")
}

// handleManager sets up the routes for the manager endpoints.
func (p *Client) handleManager(private *echo.Group, key string) {
	// Configure middleware with the custom claims type
	manager := private.Group("/manager")
	manager.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			// マネージャー用のカスタムクレームを使用
			return new(models.Claims)
		},
		SigningKey: []byte(key),
	}))
	// -H "Authorization: Bearer <token>"を付与してリクエスト
	manager.GET("/health", privateHealth)
	// 店舗の追加
	// - 登録店舗情報を取得
	manager.GET("/store", p.GetAllStores)
	// - 店舗情報を登録
	manager.POST("/store", p.RegisterStore)

}

// handleSession sets up the routes for the session endpoints.
func (p *Client) handleSession(private *echo.Group, key string) {
	// Configure middleware with the custom claims type
	session := private.Group("/session")
	session.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			// セッション用のカスタムクレームを使用
			return new(models.SessionClaims)
		},
		SigningKey: []byte(key),
	}))
	// -H "Authorization: Bearer <session_jwt>"を付与してリクエスト
	session.GET("/health", privateHealth)
}
