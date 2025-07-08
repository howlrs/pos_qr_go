package routes

import (
	"backend/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterStore は、店舗を追加するためのエンドポイントです。
func (p *Client) RegisterStore(c echo.Context) error {
	// usecases.Storeを使用して、店舗情報を受け取る
	// usecases.Storeはusecasesにてmodels.Storeに変換される
	store := &usecases.Store{}
	if err := c.Bind(store); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind store data: %v", err)
	}

	ctx := c.Request().Context()
	if err := store.Register(ctx, p.firestore); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to add store: %v", err)
	}

	return responseHandler(c, http.StatusOK, nil, nil, "Store added successfully")
}

// GetAllStores は、全ての店舗情報を取得するためのエンドポイントです。
func (p *Client) GetAllStores(c echo.Context) error {
	ctx := c.Request().Context()
	stores, err := usecases.NewStore().GetAll(ctx, p.firestore)
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to retrieve stores: %v", err)
	}

	return responseHandler(c, http.StatusOK, stores, nil, "")
}

// UpdateStore は、店舗情報を更新するためのエンドポイントです。
func (p *Client) UpdateStore(c echo.Context) error {
	// usecases.Storeを使用して、更新する店舗情報を受け取る
	store := &usecases.Store{}
	if err := c.Bind(store); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind store data: %v", err)
	}

	// リクエストから店舗IDを取得
	storeID := c.Param("id")
	if storeID == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Store ID is required")
	}

	ctx := c.Request().Context()
	if err := store.Update(ctx, p.firestore, storeID); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to update store: %v", err)
	}

	return responseHandler(c, http.StatusOK, nil, nil, "Store updated successfully")
}

// RegisterStore は、店舗を追加するためのエンドポイントです。
// DeleteStore は、店舗を削除するためのエンドポイントです。
func (p *Client) DeleteStore(c echo.Context) error {
	// リクエストから店舗IDを取得
	storeID := c.Param("id")
	if storeID == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Store ID is required")
	}

	ctx := c.Request().Context()
	if err := usecases.NewStore().Delete(ctx, p.firestore, storeID); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to delete store: %v", err)
	}

	return responseHandler(c, http.StatusOK, nil, nil, "Store deleted successfully")
}
