package mocks

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
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
