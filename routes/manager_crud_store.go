package routes

import (
	"backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestStore struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

func (s *RequestStore) IsValidate() error {
	var missingFields []string

	if s.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if s.Email == "" {
		missingFields = append(missingFields, "email")
	}
	if s.Password == "" {
		missingFields = append(missingFields, "password")
	}
	if s.Address == "" {
		missingFields = append(missingFields, "address")
	}
	if s.Phone == "" {
		missingFields = append(missingFields, "phone")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", missingFields)
	}
	return nil
}

type ResponseStore struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewResponseStore は、models.StoreをResponseStoreに変換します。
// パスワードは含まれません。
func NewResponseStore(store *models.Store) *ResponseStore {
	return &ResponseStore{
		ID:        store.ID,
		Name:      store.Name,
		Email:     store.Email,
		Address:   store.Address,
		Phone:     store.Phone,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}

// RegisterStore は、店舗を追加するためのエンドポイントです。
func (p *Client) RegisterStore(c echo.Context) error {
	// usecases.Storeを使用して、店舗情報を受け取る
	// usecases.Storeはusecasesにてmodels.Storeに変換される
	store := &RequestStore{}
	if err := c.Bind(store); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind store data: %v", err)
	}

	if err := store.IsValidate(); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Validation failed: %v", err)
	}

	createStore, err := p.uc.RegisterStore(c.Request().Context(), store.Name, store.Email, store.Password, store.Address, store.Phone)
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to register store: %v", err)
	}

	return responseHandler(c, http.StatusOK, NewResponseStore(createStore), nil, "Store added successfully")
}

// GetStore は、指定IDの店舗情報を取得するためのエンドポイントです。
func (p *Client) GetStore(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Store ID is required")
	}

	// idがあれば、特定の店舗情報を取得
	// それ以外は全ての店舗情報を取得
	store, err := p.uc.GetStore(c.Request().Context(), id)
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get store: %v", err)
	}

	return responseHandler(c, http.StatusOK, NewResponseStore(store), nil, "")
}

// GetStores は、全ての店舗情報を取得するためのエンドポイントです。
func (p *Client) GetAllStores(c echo.Context) error {
	// idがあれば、特定の店舗情報を取得
	// それ以外は全ての店舗情報を取得
	stores, err := p.uc.GetAllStores(c.Request().Context())
	if err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to get stores: %v", err)
	}

	// 返り値を整形する
	responseStores := make([]*ResponseStore, len(stores))
	for i, store := range stores {
		responseStores[i] = NewResponseStore(store)
	}

	return responseHandler(c, http.StatusOK, responseStores, nil, "")
}

// UpdateStore は、店舗情報を更新するためのエンドポイントです。
func (p *Client) UpdateStore(c echo.Context) error {
	// usecases.Storeを使用して、更新する店舗情報を受け取る
	store := &RequestStore{}
	if err := c.Bind(store); err != nil {
		return responseHandler(c, http.StatusBadRequest, nil, err, "Failed to bind store data: %v", err)
	}

	// リクエストから店舗IDを取得
	id := c.Param("id")
	if id == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Store ID is required")
	}

	if err := p.uc.Update(c.Request().Context(), id, store.Name, store.Email, store.Password, store.Address, store.Phone); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to update store: %v", err)
	}

	return responseHandler(c, http.StatusOK, nil, nil, "Store updated successfully")
}

// DeleteStore は、店舗を削除するためのエンドポイントです。
func (p *Client) DeleteStore(c echo.Context) error {
	// リクエストから店舗IDを取得
	id := c.Param("id")
	if id == "" {
		return responseHandler(c, http.StatusBadRequest, nil, nil, "Store ID is required")
	}

	if err := p.uc.Delete(c.Request().Context(), id); err != nil {
		return responseHandler(c, http.StatusInternalServerError, nil, err, "Failed to delete store: %v", err)
	}

	return responseHandler(c, http.StatusOK, nil, nil, "Store deleted successfully")
}
