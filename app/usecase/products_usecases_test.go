package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/ports/custom_error"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetProductsUseCase_List(t *testing.T) {
	tests := []struct {
		name          string
		query         input.QueryData
		repositoryErr error
		expected      output.GetAllProductsResponse
		expectedError error
		setup         func(*mocks.ProductsRepository, *mocks.ProductsPresenter, []models.Product)
	}{
		{
			name:          "rejects limit below one",
			query:         input.QueryData{Limit: 0},
			expectedError: custom_error.ErrInvalidProductLimit,
		},
		{
			name:          "rejects limit above one hundred",
			query:         input.QueryData{Limit: 101},
			expectedError: custom_error.ErrInvalidProductLimit,
		},
		{
			name:  "returns presented products with filters",
			query: input.QueryData{Offset: 10, Limit: 20, Category: "Shoes", PriceLessThan: decimal.RequireFromString("100.00")},
			expected: output.GetAllProductsResponse{
				Products: []output.Product{{Code: "PROD001", Price: 10.99}},
			},
			setup: func(repo *mocks.ProductsRepository, presenter *mocks.ProductsPresenter, products []models.Product) {
				repo.On("List", mock.Anything, input.QueryData{Offset: 10, Limit: 20, Category: "Shoes", PriceLessThan: decimal.RequireFromString("100.00")}).Return(products, nil).Once()
				presenter.On("ProductListResponse", mock.Anything, products).Return(output.GetAllProductsResponse{Products: []output.Product{{Code: "PROD001", Price: 10.99}}}).Once()
			},
		},
		{
			name:          "returns repository error",
			query:         input.QueryData{Limit: 10},
			repositoryErr: errors.New("database error"),
			expectedError: custom_error.ErrProductsFetch,
			setup: func(repo *mocks.ProductsRepository, _ *mocks.ProductsPresenter, products []models.Product) {
				repo.On("List", mock.Anything, mock.Anything).Return(products, errors.New("database error")).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.ProductsRepository)
			presenter := new(mocks.ProductsPresenter)
			products := []models.Product{{Code: "PROD001", Price: decimal.RequireFromString("10.99")}}
			if test.setup != nil {
				test.setup(repo, presenter, products)
			}

			result, err := NewProductsUseCase(repo, presenter).List(context.Background(), test.query)

			if test.expectedError != nil {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, test.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
			repo.AssertExpectations(t)
			presenter.AssertExpectations(t)
		})
	}
}

func TestGetProductsUseCase_Details(t *testing.T) {
	tests := []struct {
		name          string
		repositoryErr error
		expected      output.GetByProductCodeResponse
		expectedError error
	}{
		{
			name: "returns presented product details",
			expected: output.GetByProductCodeResponse{
				Product: output.Product{Code: "PROD001", Price: 10.99},
			},
		},
		{
			name:          "returns repository error",
			repositoryErr: gorm.ErrRecordNotFound,
			expectedError: custom_error.ErrProductNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.ProductsRepository)
			presenter := new(mocks.ProductsPresenter)
			product := models.Product{Code: "PROD001", Price: decimal.RequireFromString("10.99")}
			repo.On("Details", mock.Anything, input.QueryData{Code: "PROD001"}).Return(product, test.repositoryErr).Once()
			if test.repositoryErr == nil {
				presenter.On("ProductDetailsResponse", mock.Anything, product).Return(test.expected).Once()
			}

			result, err := NewProductsUseCase(repo, presenter).Details(context.Background(), input.QueryData{Code: "PROD001"})

			if test.expectedError != nil {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, test.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
			repo.AssertExpectations(t)
			presenter.AssertExpectations(t)
		})
	}
}
