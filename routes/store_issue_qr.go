package routes

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// IssueSeatQRForStore
// IssueSeatQRForStoreは、特定の店舗と座席のQRコード発行を処理します。
// リクエストパラメータから店舗IDと座席IDを取得し、店舗と座席を指すURLをエンコードしたQRコードを生成します。
// QRコードはbase64エンコードされた文字列としてURLとともにレスポンスで返されます。
// パラメータ抽出やQRコード生成時にエラーが発生した場合は、適切なエラーレスポンスを返します。
//
// 期待されるリクエストパラメータ:
//   - storeID: 店舗ID
//   - seatID: 座席ID
//
// レスポンス:
//   - message: QRコードが発行された店舗と座席を示す成功メッセージ
//   - url: QRコードにエンコードされたURL
//   - qr_code: base64エンコードされたQRコード画像
//
// パラメータ不正の場合はHTTP 400、内部エラーの場合はHTTP 500を返します。
func (p *Client) IssueSeatQRForStore(c echo.Context) error {
	// QRコード読み込み時にパラメータ有りウェブサイトにアクセス
	storeID, seatID, err := getStoreAndSeatID(c)
	if err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Invalid parameters: %v", err)
	}

	// QRコードの作成
	domain := os.Getenv("FRONTEND_URL")
	url := fmt.Sprintf("%s/store/%s/seat/%s", domain, storeID, seatID)
	qrCodeOfURL, err := qrcode.New(url)
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to create QR code: %v", err)
	}

	// get bytes
	buf := bytes.NewBuffer(nil)
	wc := Closer{Writer: buf}
	w2 := standard.NewWithWriter(wc, standard.WithQRWidth(40))
	if err = qrCodeOfURL.Save(w2); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to save QR code: %v", err)
	}
	defer wc.Close()

	// レスポンスを返す
	return responseHandler(c, http.StatusOK, map[string]string{
		"url":     url,
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil, fmt.Sprintf("QR code for store %s and seat %s issued successfully", storeID, seatID))
}

type Closer struct {
	io.Writer
}

func (Closer) Close() error { return nil }
