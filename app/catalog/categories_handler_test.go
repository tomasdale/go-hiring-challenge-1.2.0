package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/ports/custom_error"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/app/usecase"
	"github.com/mytheresa/go-hiring-challenge/infrastructure/logger"
	"github.com/mytheresa/go-hiring-challenge/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoriesHandler_GetAllCategories(t *testing.T) {
	t.Run("returns categories", func(t *testing.T) {
		repo := new(mocks.CategoriesRepository)
		categories := []models.Category{
			{Code: "CAT-001", Name: "Shoes"},
			{Code: "CAT-002", Name: "Bags"},
		}
		repo.On("List", mock.Anything, mock.Anything).Return(categories, nil).Once()

		presenter := new(mocks.CategoriesPresenter)
		presenter.On("CategoryListResponse", mock.Anything, categories).Return(output.ListCategoriesResponse{
			Categories: []output.Category{
				{Code: "CAT-001", Name: "Shoes"},
				{Code: "CAT-002", Name: "Bags"},
			},
		}).Once()

		h := NewCategoriesHandler(usecase.NewCategoriesUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		res := httptest.NewRecorder()

		h.GetAllCategories(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"categories":[{"code":"CAT-001","name":"Shoes"},{"code":"CAT-002","name":"Bags"}]}`, res.Body.String())
		repo.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := new(mocks.CategoriesRepository)
		repo.On("List", mock.Anything, mock.Anything).Return([]models.Category(nil), errors.New("repo error")).Once()

		presenter := new(mocks.CategoriesPresenter)
		h := NewCategoriesHandler(usecase.NewCategoriesUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		res := httptest.NewRecorder()

		h.GetAllCategories(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, `{"error":"Error fetching categories"}`, res.Body.String())
		repo.AssertExpectations(t)
		presenter.AssertNotCalled(t, "CategoryListResponse", mock.Anything, mock.Anything)
	})
}

func TestCategoriesHandler_GetCategoryByCode(t *testing.T) {
	t.Run("returns bad request when code is empty", func(t *testing.T) {
		h := NewCategoriesHandler(
			usecase.NewCategoriesUseCase(new(mocks.CategoriesRepository), new(mocks.CategoriesPresenter)),
			logger.NewNopLogger(),
		)
		req := httptest.NewRequest(http.MethodGet, "/categories/", nil)
		res := httptest.NewRecorder()

		h.GetCategoryByCode(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{"error":"category code is required"}`, res.Body.String())
	})

	repo := new(mocks.CategoriesRepository)
	repo.On("Details", mock.Anything, input.QueryData{Code: "CAT001"}).Return(models.Category{
		Code: "CAT001",
		Name: "Clothing",
	}, nil).Once()

	presenter := new(mocks.CategoriesPresenter)
	h := NewCategoriesHandler(usecase.NewCategoriesUseCase(repo, presenter), logger.NewNopLogger())
	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories/{code}", h.GetCategoryByCode)

	req := httptest.NewRequest(http.MethodGet, "/categories/CAT001", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, `{"code":"CAT001","name":"Clothing"}`, res.Body.String())
	repo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory(t *testing.T) {
	t.Run("creates category from JSON body", func(t *testing.T) {
		repo := new(mocks.CategoriesRepository)
		category := models.Category{Code: "CAT-001", Name: "Shoes"}
		repo.On("Create", mock.Anything, category).Return(nil).Once()

		h := NewCategoriesHandler(usecase.NewCategoriesUseCase(repo, new(mocks.CategoriesPresenter)), logger.NewNopLogger())
		body, err := json.Marshal(input.Category{Code: category.Code, Name: category.Name})
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		res := httptest.NewRecorder()

		h.CreateCategory(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"message":"category created successfully"}`, res.Body.String())
		repo.AssertExpectations(t)
	})

	t.Run("returns bad request for invalid JSON", func(t *testing.T) {
		repo := new(mocks.CategoriesRepository)
		h := NewCategoriesHandler(usecase.NewCategoriesUseCase(repo, new(mocks.CategoriesPresenter)), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString("invalid json"))
		res := httptest.NewRecorder()

		h.CreateCategory(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.JSONEq(t, `{"error":"invalid request body"}`, res.Body.String())
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("returns conflict when category code already exists", func(t *testing.T) {
		repo := new(mocks.CategoriesRepository)
		repo.On("Create", mock.Anything, mock.Anything).Return(custom_error.ErrCategoryCodeAlreadyExists).Once()

		h := NewCategoriesHandler(usecase.NewCategoriesUseCase(repo, new(mocks.CategoriesPresenter)), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(`{"code":"CAT-001","name":"Shoes"}`))
		res := httptest.NewRecorder()

		h.CreateCategory(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
		assert.JSONEq(t, `{"error":"category code already exists"}`, res.Body.String())
		repo.AssertExpectations(t)
	})
}
