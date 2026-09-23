package presenter

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/domain/output"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProductsPresenter_ProductListResponse(t *testing.T) {
	presenter := NewProductsPresenter()

	products := []models.Product{{
		Code:  "PRD-001",
		Price: decimal.RequireFromString("149.99"),
		Category: models.Category{
			Code: "CAT-1",
			Name: "Shoes",
		},
		Variants: []models.Variant{
			{Name: "Small", SKU: "SKU-1", Price: decimal.RequireFromString("129.99")},
			{Name: "Large", SKU: "SKU-2", Price: decimal.RequireFromString("159.99")},
		},
	}}

	expected := output.GetAllProductsResponse{
		Products: []output.Product{{
			Code:  "PRD-001",
			Price: 149.99,
			Category: output.Category{
				Code: "CAT-1",
				Name: "Shoes",
			},
			Variants: []output.Variant{
				{Name: "Small", SKU: "SKU-1", Price: 129.99},
				{Name: "Large", SKU: "SKU-2", Price: 159.99},
			},
		}},
	}

	assert.Equal(t, expected, presenter.ProductListResponse(context.Background(), products))
}

func TestProductsPresenter_ProductDetailsResponse(t *testing.T) {
	presenter := NewProductsPresenter()

	product := models.Product{
		Code:  "PRD-002",
		Price: decimal.RequireFromString("250.00"),
		Category: models.Category{
			Code: "CAT-2",
			Name: "Bags",
		},
		Variants: []models.Variant{
			{Name: "Black", SKU: "SKU-10", Price: decimal.RequireFromString("270.00")},
		},
	}

	expected := output.GetByProductCodeResponse{
		Product: output.Product{
			Code:  "PRD-002",
			Price: 250.00,
			Category: output.Category{
				Code: "CAT-2",
				Name: "Bags",
			},
			Variants: []output.Variant{{
				Name:  "Black",
				SKU:   "SKU-10",
				Price: 270.00,
			}},
		},
	}

	assert.Equal(t, expected, presenter.ProductDetailsResponse(context.Background(), product))
}
