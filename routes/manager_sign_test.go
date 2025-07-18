package routes_test

import (
	"backend/repositories"
	"backend/routes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestSinginAndPrivateRequest は、サインインAPIおよび認証が必要なプライベートAPIのテストを行います。
//
// このテストでは、まずテストユーザーでサインインし、JWTトークンを取得します。
// その後、取得したトークンを用いてプライベートエンドポイントへのリクエストを行い、認証が正しく機能しているかを検証します。
//
// 引数:
//
//	t *testing.T : テストランナーから提供されるテスト用構造体
//
// 戻り値:
//
//	なし
func TestSinginAndPrivateRequest(t *testing.T) {
	godotenv.Load("../.env.local")
	repositories.LoadConfig()

	params := []map[string]string{
		// 新規作成
		// {
		// 	"name":         "test_signup",
		// 	"url":          "/api/v1/public/signup",
		// 	"method":       "POST",
		// 	"body":         `{"email": "test_signup", "password": "test_signup"}`,
		// 	"expactStatus": "200",
		// },
		{
			"name":         "test_signin",
			"url":          "/api/v1/public/signin",
			"method":       "POST",
			"body":         `{"email": "test_signup", "password": "test_signup"}`,
			"expactStatus": "200",
		},
		{
			"name":         "test_private_health",
			"url":          "/api/v1/private/manager/health",
			"method":       "GET",
			"body":         "",
			"expactStatus": "200",
		},
	}

	// echo モックアップサーバ
	e := echo.New()
	isTest := true
	routes.Endpoint(e, isTest)

	var jwtToken string

	for _, param := range params {
		// create request
		var req *http.Request
		if jwtToken != "" {
			req = httptest.NewRequest(param["method"], param["url"], strings.NewReader(param["body"]))
			req.Header.Set("Authorization", "Bearer "+jwtToken)
		} else if param["body"] != "" {
			req = httptest.NewRequest(param["method"], param["url"], strings.NewReader(param["body"]))
		} else {
			req = httptest.NewRequest(param["method"], param["url"], nil)
		}

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var result map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err, param["name"])

		data, isThere := result["data"].(map[string]interface{})
		if isThere {
			v, ok := data["token"]
			if ok {
				jwtToken = v.(string)
			}
		}

		// status code
		if param["expactStatus"] != "" {
			// messageを取得
			message, ok := result["message"].(string)
			assert.True(t, ok)

			err, ok := result["error"].(string)
			if ok && err != "" {
				t.Errorf("error should be empty, but got %s", err)
			}

			assert.Equal(t, param["expactStatus"], fmt.Sprintf("%v", rec.Code), fmt.Sprintf("%s status code should be %s, but got %d, message: %s", param["name"], param["expactStatus"], rec.Code, message))
		}
	}
}
