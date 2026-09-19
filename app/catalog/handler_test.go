package catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/app/usecase"
	"github.com/mytheresa/go-hiring-challenge/infrastructure/logger"
	"github.com/mytheresa/go-hiring-challenge/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNumberWithDefault(t *testing.T) {
	t.Run("returns default value when string is empty", func(t *testing.T) {
		assert.Equal(t, 10, numberWithDefault("", 10))
	})

	t.Run("returns default value when conversion fails", func(t *testing.T) {
		assert.Equal(t, 7, numberWithDefault("invalid", 7))
	})

	t.Run("returns parsed value when valid", func(t *testing.T) {
		assert.Equal(t, 25, numberWithDefault("25", 10))
	})
}

func TestCatalogHandler_GetAllProducts(t *testing.T) {
	t.Run("uses default offset and limit when query is empty", func(t *testing.T) {
		repo := new(mocks.ProductsRepository)
		products := []models.Product{{Code: "PRD-001", Price: decimal.RequireFromString("99.90")}}
		expectedResponse := output.GetAllProductsResponse{
			Products: []output.ProductSummary{{Code: "PRD-001", Price: 99.9}},
		}

		repo.On("List", mock.Anything, mock.MatchedBy(func(query input.QueryData) bool {
			return query.Offset == DEFAULT_OFFSET && query.Limit == DEFAULT_LIMIT
		})).Return(products, nil).Once()

		presenter := new(mocks.ProductsPresenter)
		presenter.On("ProductListResponse", mock.Anything, products).Return(expectedResponse).Once()

		h := NewCatalogHandler(usecase.NewProductsUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		res := httptest.NewRecorder()

		h.GetAllProducts(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"products":[{"code":"PRD-001","price":99.9}]}`, res.Body.String())
		repo.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})

	t.Run("uses provided offset and limit values from query string", func(t *testing.T) {
		repo := new(mocks.ProductsRepository)
		repo.On("List", mock.Anything, mock.MatchedBy(func(query input.QueryData) bool {
			return query.Offset == 3 && query.Limit == 20
		})).Return([]models.Product{}, nil).Once()

		presenter := new(mocks.ProductsPresenter)
		presenter.On("ProductListResponse", mock.Anything, []models.Product{}).Return(output.GetAllProductsResponse{Products: []output.ProductSummary{}}).Once()

		h := NewCatalogHandler(usecase.NewProductsUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/products?offset=3&limit=20", nil)
		res := httptest.NewRecorder()

		h.GetAllProducts(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"products":[]}`, res.Body.String())
		repo.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})

	t.Run("returns not found when usecase fails", func(t *testing.T) {
		repo := new(mocks.ProductsRepository)
		repo.On("List", mock.Anything, mock.Anything).Return([]models.Product(nil), errors.New("repo error")).Once()

		presenter := new(mocks.ProductsPresenter)

		h := NewCatalogHandler(usecase.NewProductsUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		res := httptest.NewRecorder()

		h.GetAllProducts(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, `{"error":"Product not found"}`, res.Body.String())
		presenter.AssertNotCalled(t, "ProductListResponse", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})
}

func TestCatalogHandler_GetProductByCode(t *testing.T) {
	t.Run("sends product code to usecase and returns payload", func(t *testing.T) {
		repo := new(mocks.ProductsRepository)
		product := models.Product{
			Code:     "PRD-001",
			Price:    decimal.RequireFromString("150.00"),
			Category: models.Category{Code: "CAT-1", Name: "Shoes"},
			Variants: []models.Variant{{Name: "Small", SKU: "SKU-1", Price: decimal.RequireFromString("150.00")}},
		}
		repo.On("Details", mock.Anything, mock.MatchedBy(func(query input.QueryData) bool {
			return query.Code == "PRD-001"
		})).Return(product, nil).Once()

		presenter := new(mocks.ProductsPresenter)
		presenter.On("ProductDetailsResponse", mock.Anything, product).Return(output.GetByProductCodeResponse{
			Product: output.ProductDetails{
				Code:  product.Code,
				Price: product.Price.InexactFloat64(),
				Category: output.Category{
					Code: product.Category.Code,
					Name: product.Category.Name,
				},
				Variants: []output.Variant{{Name: "Small", SKU: "SKU-1", Price: 150.0}},
			},
		}).Once()

		h := NewCatalogHandler(usecase.NewProductsUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/products/PRD-001", nil)
		res := httptest.NewRecorder()

		h.GetProductByCode(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, `{"product":{"code":"PRD-001","price":150,"category":{"code":"CAT-1","name":"Shoes"},"variants":[{"name":"Small","sku":"SKU-1","price":150}]}}`, res.Body.String())
		repo.AssertExpectations(t)
		presenter.AssertExpectations(t)
	})

	t.Run("returns not found when details query fails", func(t *testing.T) {
		repo := new(mocks.ProductsRepository)
		repo.On("Details", mock.Anything, mock.Anything).Return(models.Product{}, errors.New("product not found")).Once()

		presenter := new(mocks.ProductsPresenter)

		h := NewCatalogHandler(usecase.NewProductsUseCase(repo, presenter), logger.NewNopLogger())
		req := httptest.NewRequest(http.MethodGet, "/products/UNKNOWN", nil)
		res := httptest.NewRecorder()

		h.GetProductByCode(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.JSONEq(t, `{"error":"Product not found"}`, res.Body.String())
		presenter.AssertNotCalled(t, "ProductDetailsResponse", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})
}
