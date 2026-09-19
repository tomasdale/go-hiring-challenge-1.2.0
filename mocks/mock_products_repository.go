package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type ProductsRepository struct {
	mock.Mock
}

func (m *ProductsRepository) List(ctx context.Context, query input.QueryData) ([]models.Product, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductsRepository) Details(ctx context.Context, query input.QueryData) (models.Product, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return models.Product{}, args.Error(1)
	}
	return args.Get(0).(models.Product), args.Error(1)
}

type ProductsPresenter struct {
	mock.Mock
}

func (m *ProductsPresenter) ProductListResponse(ctx context.Context, products []models.Product) output.GetAllProductsResponse {
	args := m.Called(ctx, products)
	if args.Get(0) == nil {
		return output.GetAllProductsResponse{}
	}
	return args.Get(0).(output.GetAllProductsResponse)
}

func (m *ProductsPresenter) ProductDetailsResponse(ctx context.Context, product models.Product) output.GetByProductCodeResponse {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return output.GetByProductCodeResponse{}
	}
	return args.Get(0).(output.GetByProductCodeResponse)
}
