package custom_error

import "errors"

var ErrCategoryCodeAlreadyExists = errors.New("category code already exists")

var ErrProductCodeAlreadyExists = errors.New("product code already exists")

var ErrInvalidProductLimit = errors.New("Limit must be between 1 and 100")

var ErrProductsFetch = errors.New("Error fetching products")

var ErrProductNotFound = errors.New("product not found")

var ErrProductDetailsFetch = errors.New("Error fetching product details")

var ErrCategoriesFetch = errors.New("Error fetching categories")

var ErrCategoryNotFound = errors.New("category not found")

var ErrCategoryDetailsFetch = errors.New("Error fetching category details")

var ErrCategoryInputInvalid = errors.New("Category code and name cannot be empty")

var ErrCategoryCreation = errors.New("error creating category")
