package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mytheresa/go-hiring-challenge/app/ports/custom_error"
	"github.com/mytheresa/go-hiring-challenge/app/ports/input"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) List(ctx context.Context, query input.QueryData) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) Details(ctx context.Context, query input.QueryData) (models.Category, error) {
	if query.Code == "" {
		return models.Category{}, gorm.ErrRecordNotFound
	}

	var category models.Category
	if err := r.db.WithContext(ctx).Where("code = ?", query.Code).First(&category).Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (r *CategoriesRepository) Create(ctx context.Context, category models.Category) error {
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) && postgresError.Code == "23505" {
			return custom_error.ErrCategoryCodeAlreadyExists
		}
		return err
	}
	return nil
}
