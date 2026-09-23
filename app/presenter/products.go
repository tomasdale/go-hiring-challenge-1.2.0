package presenter

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/domain/output"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductsPresenter struct {
}

func NewProductsPresenter() *ProductsPresenter {
	return &ProductsPresenter{}
}

func (d *ProductsPresenter) ProductListResponse(ctx context.Context, i []models.Product) output.GetAllProductsResponse {
	input := i

	products := make([]output.Product, len(input))
	for i, p := range input {
		products[i] = output.Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: output.Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
			Variants: func() []output.Variant {
				variants := make([]output.Variant, len(p.Variants))
				for i, v := range p.Variants {
					variants[i] = output.Variant{
						Name:  v.Name,
						SKU:   v.SKU,
						Price: v.Price.InexactFloat64(),
					}
				}
				return variants
			}(),
		}
	}

	response := output.GetAllProductsResponse{
		Products: products,
	}
	return response
}

func (d *ProductsPresenter) ProductDetailsResponse(ctx context.Context, product models.Product) output.GetByProductCodeResponse {
	variants := make([]output.Variant, len(product.Variants))
	for i, variant := range product.Variants {
		variants[i] = output.Variant{
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: variant.Price.InexactFloat64(),
		}
	}

	return output.GetByProductCodeResponse{
		Product: output.Product{
			Code:  product.Code,
			Price: product.Price.InexactFloat64(),
			Category: output.Category{
				Code: product.Category.Code,
				Name: product.Category.Name,
			},
			Variants: variants,
		},
	}
}
