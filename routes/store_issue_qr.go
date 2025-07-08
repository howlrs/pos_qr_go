package routes

import (
	"backend/usecases"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (p *Client) IssueSeatQRForStore(c echo.Context) error {
	// QRコード読み込み時にパラメータ有りウェブサイトにアクセス
	storeID, seatID, err := getStoreAndSeatID(c)
	if err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Invalid parameters: %v", err)
	}

	// セッションJWT発行用の構造体を生成
	session := usecases.NewSession(storeID, seatID, time.Now().UTC().Add(1*time.Hour))
	token, err := session.CreateJWT()
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to create session JWT: %v", err)
	}

	// レスポンスを返す
	return responseHandler(c, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Session started for store_id=%s, seat_id=%s", storeID, seatID),
		"token":   token,
	}, nil, "")
}
