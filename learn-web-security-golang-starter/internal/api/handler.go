package api

import (
	"net/http"
	"strings"

	"github.com/bootdotdev/learn-web-security/internal/accounts"
	"github.com/bootdotdev/learn-web-security/internal/auth/sessions"
	"github.com/bootdotdev/learn-web-security/internal/httpx"
	"github.com/bootdotdev/learn-web-security/internal/logging"
	"github.com/bootdotdev/learn-web-security/internal/orders"
	"github.com/bootdotdev/learn-web-security/internal/storefront"
)

type integrationOrderResponse struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`
	TotalCents int64  `json:"total_cents"`
	CreatedAt  string `json:"created_at"`
}

type orderItemResponse struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
	PriceCents  int64  `json:"price_cents"`
}

func orderItemToOrderItemResponse(oi orders.Item) orderItemResponse {
	return orderItemResponse{
		ProductID:   oi.ProductID,
		ProductName: oi.ProductName,
		Quantity:    oi.Quantity,
		PriceCents:  oi.PriceCents,
	}
}

type productResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImagePath   string `json:"image_path"`
	PriceCents  int64  `json:"price_cents"`
}

func productToProductResponse(p storefront.Product) productResponse {
	return productResponse{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ImagePath:   p.ImagePath,
		PriceCents:  p.PriceCents,
	}
}

type orderResponse struct {
	Id         int64  `json:"id"`
	Status     string `json:"status"`
	TotalCents int64  `json:"total_cents"`
	CreatedAt  string `json:"created_at"`
}

func orderToOrderResponse(o orders.Order) orderResponse {
	return orderResponse{
		Id:         o.ID,
		Status:     o.Status,
		TotalCents: o.TotalCents,
		CreatedAt:  o.CreatedAt,
	}
}

type Handler struct {
	accountStore      *accounts.Store
	orderStore        *orders.Store
	productStore      *storefront.Store
	apiStore          *Store
	logger            *logging.Logger
	maxProductResults int64
}

func NewHandler(accountStore *accounts.Store, orderStore *orders.Store, productStore *storefront.Store, apiStore *Store, logger *logging.Logger, maxProductResults int) *Handler {
	return &Handler{
		accountStore: accountStore, orderStore: orderStore, productStore: productStore, apiStore: apiStore,
		logger: logger, maxProductResults: int64(maxProductResults),
	}
}

func (handler *Handler) AccountOrders(responseWriter http.ResponseWriter, request *http.Request) {
	current, ok := handler.requireAuthentication(responseWriter, request)
	if !ok {
		return
	}
	orders, err := handler.orderStore.ListForUser(request.Context(), current.User.ID)
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return
	}
	ordersResponse := make([]orderResponse, len(orders))
	for i, o := range orders {
		ordersResponse[i] = orderToOrderResponse(o)
	}
	httpx.RespondWithJSON(responseWriter, http.StatusOK, map[string]any{"orders": ordersResponse})
}

func (handler *Handler) Order(responseWriter http.ResponseWriter, request *http.Request) {
	current, ok := handler.requireAuthentication(responseWriter, request)
	if !ok {
		return
	}
	orderID, valid := httpx.ParseSafeInteger(request.PathValue("id"))
	if !valid {
		httpx.RespondWithJSON(responseWriter, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	}
	order, found, err := handler.orderStore.FindByID(request.Context(), orderID)
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return
	}
	if !found || order.UserID != current.User.ID {
		httpx.RespondWithJSON(responseWriter, http.StatusNotFound, map[string]string{"error": "Order not found"})
		return
	}
	items, err := handler.orderStore.ListItems(request.Context(), order.ID)
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return
	}
	itemResponses := make([]orderItemResponse, 0, len(items))
	for _, item := range items {
		itemResponses = append(itemResponses, orderItemToOrderItemResponse(item))
	}
	httpx.RespondWithJSON(responseWriter, http.StatusOK, map[string]any{"order": orderToOrderResponse(order), "items": itemResponses})
}

func (handler *Handler) Products(responseWriter http.ResponseWriter, request *http.Request) {
	products, err := handler.productStore.ListProducts(request.Context(), handler.maxProductResults)
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return
	}
	var productsResponse []productResponse
	for _, p := range products {
		productsResponse = append(productsResponse, productToProductResponse(p))
	}

	httpx.RespondWithJSON(responseWriter, http.StatusOK, map[string]any{"products": productsResponse})
}

func (handler *Handler) WarehouseOrders(responseWriter http.ResponseWriter, request *http.Request) {
	apiKey := request.Header.Get("x-api-key")
	key, found, err := handler.apiStore.FindKey(request.Context(), apiKey)
	if err != nil || !found {
		msg := "Not allowed"
		if err != nil {
			msg = err.Error()
		}
		httpx.RespondWithError(responseWriter, http.StatusUnauthorized, msg)
		return
	}
	if !strings.Contains(key.Scope, "orders:read") {
		httpx.RespondWithError(responseWriter, http.StatusForbidden, "Not allowed")
		return
	}
	quota, err := handler.apiStore.ConsumeQuota(request.Context(), key.ID)
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return
	}
	SetQuotaHeaders(responseWriter, quota)
	if !quota.Allowed {
		RespondWithQuotaExhausted(responseWriter, quota)
		return
	}
	orders, err := handler.orderStore.ListAll(request.Context())
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return
	}
	responses := make([]integrationOrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, integrationOrderResponse{
			ID: order.ID, Status: order.Status, TotalCents: order.TotalCents, CreatedAt: order.CreatedAt,
		})
	}
	httpx.RespondWithJSON(responseWriter, http.StatusOK, map[string]any{
		"integration": "Warehouse Fulfillment Integration",
		"orders":      responses,
		"quota":       ToQuotaResponse(quota),
	})
}

func (handler *Handler) requireAuthentication(responseWriter http.ResponseWriter, request *http.Request) (accounts.CurrentSession, bool) {
	current, found, err := sessions.Current(request, handler.accountStore)
	if err != nil {
		handler.internalError(responseWriter, request, err)
		return accounts.CurrentSession{}, false
	}
	if !found {
		httpx.RespondWithJSON(responseWriter, http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
		return accounts.CurrentSession{}, false
	}
	return current, true
}

func (handler *Handler) internalError(responseWriter http.ResponseWriter, request *http.Request, err error) {
	_ = handler.logger.Event("unhandled_error", map[string]any{"method": request.Method, "path": request.URL.Path, "message": err.Error()})
	httpx.RespondWithError(responseWriter, http.StatusInternalServerError, err.Error())
}
