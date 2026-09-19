package ports

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/ports/output"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductsPresenter interface {
	ProductListResponse(ctx context.Context, i []models.Product) output.GetAllProductsResponse
	ProductDetailsResponse(ctx context.Context, product models.Product) output.GetByProductCodeResponse
}
