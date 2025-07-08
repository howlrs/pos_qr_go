package routes

import (
	"backend/usecases"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// StartSession は、席ユーザーのために新しいセッションを開始します。
// QRコード読み込み時に呼び出され、store_idとseat_idをクエリパラメータから取得し、
// セッション用JWTトークンを生成してCookieにセットします。
// 有効期限はデフォルトで1時間後ですが、"exp"コンテキスト値が存在する場合はその値を使用します。
//
// 引数:
//
//	c echo.Context : Echoのコンテキスト。リクエスト情報やレスポンス操作に使用します。
//
// 戻り値:
//
//	error : エラーが発生した場合はその内容を返します。正常時はnilを返します。
func (p *Client) StartSession(c echo.Context) error {
	// 着座（QR読み込み）時にパラメータ有りウェブサイトにアクセス
	storeID, seatID, err := getStoreAndSeatID(c)
	if err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Invalid parameters: %v", err)
	}

	// クエリパラメータから有効期限を取得
	// 設定がなければ、デフォルトで1時間後に設定
	expiredAt := time.Now().UTC().Add(1 * time.Hour) // デフォルトの有効期限
	exp := c.Get("exp").(int64)
	// 設定値があれば設定値を代入
	if exp > 0 {
		expiredAt = time.Unix(exp, 0).UTC() // クエリパラメータから取得した有効期限を使用
	}

	// セッションJWT発行用の構造体を生成
	// JWTを発行
	session := usecases.NewSession(storeID, seatID, expiredAt)
	token, err := session.CreateJWT()
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to create session JWT: %v", err)
	}

	// Set-Cookiesにして返す
	// 以後、当Tokenを使用して、セッションの有効性を確認する
	// 当Cookieの有効期限を設定し、ブラウザで自動的にCookieを削除
	// 当処理で、ブラウザにてSession JWTがなければ再度QRコードを読み込むことを促す
	c.SetCookie(&http.Cookie{
		Name:     "session_jwt",
		Value:    token,
		Expires:  expiredAt,
		HttpOnly: true, // JavaScriptからはアクセスできないようにする
		Secure:   true, // HTTPSでのみ送信されるようにする
	})

	// レスポンスを返す
	return responseHandler(c, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Session started for store_id=%s, seat_id=%s", storeID, seatID),
		"token":   token,
	}, nil, "Session started for store_id=%s, seat_id=%s", storeID, seatID)
}

// getStoreAndSeatID は、クエリパラメータから必須情報であるstore_idとseat_idを取得します。
// もしどちらかが存在しない場合は、エラーを返します。
//
// 引数:
//
//	c echo.Context : Echoのコンテキスト。
//
// 戻り値:
//
//	string : store_id
//	string : seat_id
//	error  : store_idまたはseat_idが存在しない場合はエラーを返します。
func getStoreAndSeatID(c echo.Context) (string, string, error) {
	storeID := c.QueryParam("store_id")
	seatID := c.QueryParam("seat_id")

	if storeID == "" || seatID == "" {
		return "", "", fmt.Errorf("store_id and seat_id are required parameters, but got store_id=%s, seat_id=%s", storeID, seatID)
	}

	return storeID, seatID, nil
}
