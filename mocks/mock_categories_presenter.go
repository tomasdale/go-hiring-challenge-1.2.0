package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type CategoriesPresenter struct {
	mock.Mock
}

func (m *CategoriesPresenter) CategoryListResponse(ctx context.Context, categories []models.Category) output.ListCategoriesResponse {
	args := m.Called(ctx, categories)
	if args.Get(0) == nil {
		return output.ListCategoriesResponse{}
	}
	return args.Get(0).(output.ListCategoriesResponse)
}
