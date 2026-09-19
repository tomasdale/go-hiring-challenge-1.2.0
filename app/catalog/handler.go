package catalog

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/usecase"
	"github.com/mytheresa/go-hiring-challenge/infrastructure/logger"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	DEFAULT_OFFSET = 0
	DEFAULT_LIMIT  = 10
)

type CatalogHandler struct {
	usecase *usecase.GetProductsUseCase
	logger  logger.Logger
}

func NewCatalogHandler(uc *usecase.GetProductsUseCase, lg logger.Logger) *CatalogHandler {
	if lg == nil {
		lg = logger.NewNopLogger()
	}

	return &CatalogHandler{
		usecase: uc,
		logger:  lg,
	}
}

func (h *CatalogHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	offset := numberWithDefault(query.Get("offset"), DEFAULT_OFFSET)
	limit := numberWithDefault(query.Get("limit"), DEFAULT_LIMIT)
	categoryFilter := query.Get("category")

	// Ignoring error as a request without price filter is allowed
	priceFilter, _ := decimal.NewFromString(query.Get("priceLessThan"))

	queryData := input.QueryData{
		Offset:        offset,
		Limit:         limit,
		Category:      categoryFilter,
		PriceLessThan: priceFilter,
	}

	result, err := h.usecase.List(ctx, queryData)
	if err != nil {
		h.logger.Error("failed to list products", zap.String("offset", strconv.Itoa(queryData.Offset)), zap.String("limit", strconv.Itoa(queryData.Limit)), zap.Error(err))
		api.ErrorResponse(w, err.Code(), err.Error())
		return
	}

	if len(result.Products) == 0 {
		api.ErrorResponse(w, http.StatusNotFound, "No products found")
		return
	}

	api.OKResponse(w, result)
}

func numberWithDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}

func (h *CatalogHandler) GetProductByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := r.PathValue("code")
	if code == "" {
		code = r.URL.Query().Get("code")
	}
	if code == "" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) > 1 {
			code = parts[1]
		}
	}

	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	result, err := h.usecase.Details(ctx, input.QueryData{Code: code})
	if err != nil {
		h.logger.Error("failed to get product by code", zap.String("code", code), zap.Error(err))
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	api.OKResponse(w, result)
}
