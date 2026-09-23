package usecase

import (
	"context"
	"errors"

	"github.com/mytheresa/go-hiring-challenge/app/ports"
	"github.com/mytheresa/go-hiring-challenge/app/ports/custom_error"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"gorm.io/gorm"
)

type ProductsUseCase struct {
	repo      ports.ProductsRepository
	presenter ports.ProductsPresenter
}

func NewProductsUseCase(repo ports.ProductsRepository, presenter ports.ProductsPresenter) *ProductsUseCase {
	return &ProductsUseCase{
		repo:      repo,
		presenter: presenter,
	}
}

func (uc *ProductsUseCase) List(ctx context.Context, query input.QueryData) (output.GetAllProductsResponse, error) {
	if query.Limit > 100 || query.Limit < 1 {
		return output.GetAllProductsResponse{}, custom_error.ErrInvalidProductLimit
	}

	res, err := uc.repo.List(ctx, query)
	if err != nil {
		return output.GetAllProductsResponse{}, custom_error.ErrProductsFetch
	}

	return uc.presenter.ProductListResponse(ctx, res), nil
}

func (uc *ProductsUseCase) Details(ctx context.Context, query input.QueryData) (output.GetByProductCodeResponse, error) {
	res, err := uc.repo.Details(ctx, query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return output.GetByProductCodeResponse{}, custom_error.ErrProductNotFound
		}
		return output.GetByProductCodeResponse{}, custom_error.ErrProductDetailsFetch
	}

	return uc.presenter.ProductDetailsResponse(ctx, res), nil
}
