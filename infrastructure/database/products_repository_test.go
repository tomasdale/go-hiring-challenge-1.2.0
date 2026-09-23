package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariantsWithEffectivePriceSelect(t *testing.T) {
	expected := "product_variants.id, product_variants.product_id, product_variants.name, product_variants.sku, COALESCE(product_variants.price, products.price) AS price"

	assert.Equal(t, expected, variantsWithEffectivePriceSelect())
}
