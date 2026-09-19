package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type stubProductsPresenter struct {
	listFn    func(ctx context.Context, products []models.Product) output.GetAllProductsResponse
	detailsFn func(ctx context.Context, product models.Product) output.GetByProductCodeResponse
}

func (s *stubProductsPresenter) ProductListResponse(ctx context.Context, products []models.Product) output.GetAllProductsResponse {
	if s.listFn != nil {
		return s.listFn(ctx, products)
	}
	return output.GetAllProductsResponse{}
}

func (s *stubProductsPresenter) ProductDetailsResponse(ctx context.Context, product models.Product) output.GetByProductCodeResponse {
	if s.detailsFn != nil {
		return s.detailsFn(ctx, product)
	}
	return output.GetByProductCodeResponse{}
}
