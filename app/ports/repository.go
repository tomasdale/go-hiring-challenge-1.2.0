package ports

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductsRepository interface {
	List(ctx context.Context, query input.QueryData) ([]models.Product, error)
	Details(ctx context.Context, query input.QueryData) (models.Product, error)
}

type CategoriesRepository interface {
	List(ctx context.Context, query input.QueryData) ([]models.Category, error)
	Details(ctx context.Context, query input.QueryData) (models.Category, error)
	Create(ctx context.Context, category input.Category) error
}
