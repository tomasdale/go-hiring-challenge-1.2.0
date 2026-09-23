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

type CategoriesUseCase struct {
	repo      ports.CategoriesRepository
	presenter ports.CategoriesPresenter
}

func NewCategoriesUseCase(repo ports.CategoriesRepository, presenter ports.CategoriesPresenter) *CategoriesUseCase {
	return &CategoriesUseCase{repo: repo, presenter: presenter}
}

func (c *CategoriesUseCase) List(ctx context.Context, query input.QueryData) (output.ListCategoriesResponse, error) {
	categories, err := c.repo.List(ctx, query)
	if err != nil {
		return output.ListCategoriesResponse{}, custom_error.ErrCategoriesFetch
	}

	return c.presenter.CategoryListResponse(ctx, categories), nil
}

func (c *CategoriesUseCase) Details(ctx context.Context, query input.QueryData) (output.Category, error) {
	category, err := c.repo.Details(ctx, query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return output.Category{}, custom_error.ErrCategoryNotFound
		}
		return output.Category{}, custom_error.ErrCategoryDetailsFetch
	}

	return output.Category{
		Code: category.Code,
		Name: category.Name,
	}, nil
}

func (c *CategoriesUseCase) Create(ctx context.Context, category input.Category) error {
	if category.Code == "" || category.Name == "" {
		return custom_error.ErrCategoryInputInvalid
	}

	err := c.repo.Create(ctx, category)
	if err != nil {
		if errors.Is(err, custom_error.ErrCategoryCodeAlreadyExists) {
			return custom_error.ErrCategoryCodeAlreadyExists
		}
		return custom_error.ErrCategoryCreation
	}

	return nil
}
