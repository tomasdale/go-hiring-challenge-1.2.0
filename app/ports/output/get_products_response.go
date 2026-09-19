package output

type GetAllProductsResponse struct {
	Products []ProductSummary `json:"products"`
}

type GetByProductCodeResponse struct {
	Product ProductDetails `json:"product"`
}

type ProductSummary struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type ProductDetails struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category Category  `json:"category"`
	Variants []Variant `json:"variants"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}
