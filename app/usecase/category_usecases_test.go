package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/domain/custom_error"
	"github.com/mytheresa/go-hiring-challenge/app/domain/input"
	"github.com/mytheresa/go-hiring-challenge/app/domain/output"
	"github.com/mytheresa/go-hiring-challenge/mocks"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoriesUseCase_List(t *testing.T) {
	categories := []models.Category{{Code: "CAT001", Name: "Clothing"}}
	expected := output.ListCategoriesResponse{Categories: []output.Category{{Code: "CAT001", Name: "Clothing"}}}
	tests := []struct {
		name          string
		repositoryErr error
		expected      output.ListCategoriesResponse
		expectedError error
	}{
		{name: "returns presented categories", expected: expected},
		{name: "returns repository error", repositoryErr: errors.New("database error"), expectedError: custom_error.ErrCategoriesFetch},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.CategoriesRepository)
			presenter := new(mocks.CategoriesPresenter)
			repo.On("List", mock.Anything, input.QueryData{}).Return(categories, test.repositoryErr).Once()
			if test.repositoryErr == nil {
				presenter.On("CategoryListResponse", mock.Anything, categories).Return(expected).Once()
			}

			result, err := NewCategoriesUseCase(repo, presenter).List(context.Background(), input.QueryData{})

			assertCategoryResult(t, result, err, test.expected, test.expectedError)
			repo.AssertExpectations(t)
			presenter.AssertExpectations(t)
		})
	}
}

func TestCategoriesUseCase_Details(t *testing.T) {
	category := models.Category{Code: "CAT001", Name: "Clothing"}
	tests := []struct {
		name          string
		repositoryErr error
		expected      output.Category
		expectedError error
	}{
		{name: "returns category", expected: output.Category{Code: "CAT001", Name: "Clothing"}},
		{name: "returns repository error", repositoryErr: errors.New("not found"), expectedError: custom_error.ErrCategoryDetailsFetch},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.CategoriesRepository)
			presenter := new(mocks.CategoriesPresenter)
			repo.On("Details", mock.Anything, input.QueryData{Code: "CAT001"}).Return(category, test.repositoryErr).Once()

			result, err := NewCategoriesUseCase(repo, presenter).Details(context.Background(), input.QueryData{Code: "CAT001"})

			assertCategoryResult(t, result, err, test.expected, test.expectedError)
			repo.AssertExpectations(t)
		})
	}
}

func TestCategoriesUseCase_Create(t *testing.T) {
	tests := []struct {
		name          string
		category      input.Category
		repositoryErr error
		expectedError error
	}{
		{name: "rejects empty code", category: input.Category{Name: "Clothing"}, expectedError: custom_error.ErrCategoryInputInvalid},
		{name: "rejects empty name", category: input.Category{Code: "CAT001"}, expectedError: custom_error.ErrCategoryInputInvalid},
		{name: "creates category", category: input.Category{Code: "CAT001", Name: "Clothing"}},
		{name: "returns conflict for duplicate code", category: input.Category{Code: "CAT001", Name: "Clothing"}, repositoryErr: custom_error.ErrCategoryCodeAlreadyExists, expectedError: custom_error.ErrCategoryCodeAlreadyExists},
		{name: "returns internal error", category: input.Category{Code: "CAT001", Name: "Clothing"}, repositoryErr: errors.New("database error"), expectedError: custom_error.ErrCategoryCreation},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(mocks.CategoriesRepository)
			if test.category.Code != "" && test.category.Name != "" {
				repo.On("Create", mock.Anything, input.Category{Code: test.category.Code, Name: test.category.Name}).Return(test.repositoryErr).Once()
			}

			err := NewCategoriesUseCase(repo, new(mocks.CategoriesPresenter)).Create(context.Background(), test.category)

			if test.expectedError == nil {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.ErrorIs(t, err, test.expectedError)
			}
			repo.AssertExpectations(t)
		})
	}
}

func assertCategoryResult[T any](t *testing.T, result T, err error, expected T, expectedError error) {
	t.Helper()
	if expectedError == nil {
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		return
	}
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, expectedError)
	}
}
