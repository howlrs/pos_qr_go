package routes_test

import (
	"backend/repositories"
	"backend/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIssueSeatQRForStore(t *testing.T) {
	// Load environment variables for JWT_SECRET (ファイルが存在しない場合は無視)
	_ = godotenv.Load("../.env.local")
	repositories.LoadConfig()

	// テスト用のJWT_SECRETを設定（環境変数が未設定の場合）
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test_secret_for_testing")
	}

	tests := []struct {
		name           string
		storeID        string
		seatID         string
		expectedStatus int
		expectToken    bool
		expectError    bool
	}{
		{
			name:           "正常なパラメータでのQR発行",
			storeID:        "store123",
			seatID:         "seat456",
			expectedStatus: http.StatusOK,
			expectToken:    true,
			expectError:    false,
		},
		{
			name:           "store_idが空の場合",
			storeID:        "",
			seatID:         "seat456",
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
			expectError:    true,
		},
		{
			name:           "seat_idが空の場合",
			storeID:        "store123",
			seatID:         "",
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
			expectError:    true,
		},
		{
			name:           "両方のパラメータが空の場合",
			storeID:        "",
			seatID:         "",
			expectedStatus: http.StatusBadRequest,
			expectToken:    false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { // Echo インスタンスとルートクライアントを作成
			e := echo.New()
			client := routes.NewClient(true)

			// リクエストURLを構築
			url := "/api/v1/store/qr"
			if tt.storeID != "" || tt.seatID != "" {
				url += "?"
				if tt.storeID != "" {
					url += "store_id=" + tt.storeID
				}
				if tt.seatID != "" {
					if tt.storeID != "" {
						url += "&"
					}
					url += "seat_id=" + tt.seatID
				}
			}

			// HTTPリクエストを作成
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// エンドポイントを実行
			err := client.IssueSeatQRForStore(c)

			// エラーが発生しないことを確認（Echoのエラーハンドリング）
			assert.NoError(t, err)

			// ステータスコードを確認
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// レスポンスボディを解析
			var response map[string]interface{}
			err = json.NewDecoder(rec.Body).Decode(&response)
			assert.NoError(t, err)

			if tt.expectError {
				// エラーレスポンスの確認
				assert.Contains(t, response, "error")
				assert.Contains(t, response, "message")
				assert.NotEmpty(t, response["error"])
			} else {
				// 成功レスポンスの確認
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "data")

				// dataフィールドの中身を確認
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok, "data should be a map")

				if tt.expectToken {
					assert.Contains(t, data, "qr_code")
					assert.NotEmpty(t, data["qr_code"], "qr_code should not be empty")

					// トークンが文字列であることを確認
					qrcode, ok := data["qr_code"].(string)
					assert.True(t, ok, "qr_code should be a string")
					assert.NotEmpty(t, qrcode, "qr_code should not be empty")
				}

				domain := os.Getenv("FRONTEND_URL")
				assert.Equal(t, data["url"], domain+"/store/"+tt.storeID+"/seat/"+tt.seatID)
				t.Logf("QR code URL: %s/store/%s/seat/%s", domain, tt.storeID, tt.seatID)
			}
		})
	}
}

func TestIssueSeatQRForStore_Integration(t *testing.T) {
	// Load environment variables (ファイルが存在しない場合は無視)
	_ = godotenv.Load("../.env.local")
	repositories.LoadConfig()

	// テスト用のJWT_SECRETを設定（環境変数が未設定の場合）
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test_secret_for_testing")
	}

	// 正常なケースのテスト
	t.Run("統合テスト - 正常なQR発行", func(t *testing.T) {
		e := echo.New()
		client := routes.NewClient(true)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/store/qr?store_id=test_store&seat_id=test_seat", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := client.IssueSeatQRForStore(c)
		assert.NoError(t, err)

		var response map[string]interface{}
		assert.NoError(t, json.NewDecoder(rec.Body).Decode(&response))

		assert.Contains(t, response, "message")
		assert.Contains(t, response, "data")

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok, "data should be a map")
		assert.Contains(t, data, "qr_code")
		assert.NotEmpty(t, data["qr_code"], "qr_code should not be empty")
	})

	t.Run("統合テスト - パラメータ不足", func(t *testing.T) {
		e := echo.New()
		client := routes.NewClient(true)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/store/qr", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// エンドポイントを実行
		err := client.IssueSeatQRForStore(c)
		assert.NoError(t, err)

		// ステータスコードとレスポンスの検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		assert.NoError(t, json.NewDecoder(rec.Body).Decode(&response))
		assert.Contains(t, response, "error")
		assert.Contains(t, response, "message")
	})
}
