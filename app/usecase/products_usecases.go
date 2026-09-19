package usecase

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports"
	"github.com/mytheresa/go-hiring-challenge/app/ports/custom_error"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
)

type GetProductsUseCase struct {
	repo      ports.ProductsRepository
	presenter ports.ProductsPresenter
}

func NewProductsUseCase(repo ports.ProductsRepository, presenter ports.ProductsPresenter) *GetProductsUseCase {
	return &GetProductsUseCase{
		repo:      repo,
		presenter: presenter,
	}
}

func (uc *GetProductsUseCase) List(ctx context.Context, query input.QueryData) (output.GetAllProductsResponse, ports.Error) {
	if query.Limit > 100 || query.Limit < 1 {
		return output.GetAllProductsResponse{}, custom_error.New("Limit must be between 1 and 100", 400)
	}

	res, err := uc.repo.List(ctx, query)
	if err != nil {
		return output.GetAllProductsResponse{}, custom_error.New("Error fetching products", 500)
	}

	return uc.presenter.ProductListResponse(ctx, res), nil
}

func (uc *GetProductsUseCase) Details(ctx context.Context, query input.QueryData) (output.GetByProductCodeResponse, ports.Error) {
	res, err := uc.repo.Details(ctx, query)
	if err != nil {
		return output.GetByProductCodeResponse{}, custom_error.New("Error fetching product details", 500)
	}

	return uc.presenter.ProductDetailsResponse(ctx, res), nil
}
