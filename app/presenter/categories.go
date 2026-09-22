package presenter

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesPresenter struct {
}

func NewCategoriesPresenter() *CategoriesPresenter {
	return &CategoriesPresenter{}
}

func (d *CategoriesPresenter) CategoryListResponse(ctx context.Context, input []models.Category) output.ListCategoriesResponse {
	categories := make([]output.Category, len(input))
	for i, category := range input {
		categories[i] = output.Category{
			Code: category.Code,
			Name: category.Name,
		}
	}

	return output.ListCategoriesResponse{
		Categories: categories,
	}
}
