package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/domain/output"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

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
