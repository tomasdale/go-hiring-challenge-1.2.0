package database

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/app/domain/input"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	CATEGORY_TABLE = "Category"
	VARIANTS_TABLE = "Variants"
	OVERRIDE_PRICE = "COALESCE(product_variants.price, products.price) AS price"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) List(ctx context.Context, data input.QueryData) ([]models.Product, error) {
	var query *gorm.DB

	query = r.db.WithContext(ctx).
		Model(&models.Product{}).
		Preload(VARIANTS_TABLE, variantsWithEffectivePrice()).
		Preload(CATEGORY_TABLE)
	query = withCategory(query, data.Category)
	query = withPriceLimit(query, data.PriceLessThan)

	var products []models.Product
	if err := query.Offset(data.Offset).Limit(data.Limit).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepository) Details(ctx context.Context, query input.QueryData) (models.Product, error) {
	if query.Code == "" {
		return models.Product{}, gorm.ErrRecordNotFound
	}

	var product models.Product

	if err := r.db.WithContext(ctx).
		Joins(CATEGORY_TABLE).
		Preload(VARIANTS_TABLE, variantsWithEffectivePrice()).
		Where("products.code = ?", query.Code).
		First(&product).Error; err != nil {

		return models.Product{}, err
	}

	return product, nil
}

func withCategory(query *gorm.DB, category string) *gorm.DB {
	if category == "" {
		return query
	}
	return query.Joins("JOIN categories ON categories.id = products.category_id").
		Where("categories.name = ?", category)
}

func withPriceLimit(query *gorm.DB, priceLimit decimal.Decimal) *gorm.DB {
	if priceLimit.IsZero() {
		return query
	}
	return query.Where("products.price < ?", priceLimit)
}

func variantsWithEffectivePrice() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(variantsWithEffectivePriceSelect()).
			Joins("JOIN products ON products.id = product_variants.product_id")
	}
}

func variantsWithEffectivePriceSelect() string {
	return "product_variants.id, product_variants.product_id, product_variants.name, product_variants.sku, " + OVERRIDE_PRICE
}
