package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/domain/input"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/mock"
)

type CategoriesRepository struct {
	mock.Mock
}

func (m *CategoriesRepository) List(ctx context.Context, query input.QueryData) ([]models.Category, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *CategoriesRepository) Details(ctx context.Context, query input.QueryData) (models.Category, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return models.Category{}, args.Error(1)
	}
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *CategoriesRepository) Create(ctx context.Context, category input.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}
