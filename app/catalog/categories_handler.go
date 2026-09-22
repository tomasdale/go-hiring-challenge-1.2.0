package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/usecase"
	"github.com/mytheresa/go-hiring-challenge/infrastructure/logger"
)

type CategoriesHandler struct {
	categories *usecase.CategoriesUseCase
	logger     logger.Logger
}

func NewCategoriesHandler(categories *usecase.CategoriesUseCase, l logger.Logger) *CategoriesHandler {
	return &CategoriesHandler{
		categories: categories,
		logger:     l,
	}
}

func (h *CategoriesHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.categories.List(r.Context(), input.QueryData{})
	if err != nil {
		api.ErrorResponseFromError(w, err)
		return
	}

	api.OKResponse(w, result)
}

func (h *CategoriesHandler) GetCategoryByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "category code is required")
		return
	}

	result, err := h.categories.Details(r.Context(), input.QueryData{Code: code})
	if err != nil {
		api.ErrorResponseFromError(w, err)
		return
	}

	api.OKResponse(w, result)
}

func (h *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category input.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.categories.Create(r.Context(), category)
	if err != nil {
		api.ErrorResponseFromError(w, err)
		return
	}

	api.OKResponse(w, map[string]string{"message": "category created successfully"})
}
