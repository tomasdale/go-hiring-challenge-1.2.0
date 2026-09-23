package presenter

import (
	"context"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/domain/output"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestCategoriesPresenter_CategoryListResponse(t *testing.T) {
	presenter := NewCategoriesPresenter()

	categories := []models.Category{
		{Code: "CAT-1", Name: "Shoes"},
		{Code: "CAT-2", Name: "Bags"},
	}

	expected := output.ListCategoriesResponse{
		Categories: []output.Category{
			{Code: "CAT-1", Name: "Shoes"},
			{Code: "CAT-2", Name: "Bags"},
		},
	}

	assert.Equal(t, expected, presenter.CategoryListResponse(context.Background(), categories))
}
