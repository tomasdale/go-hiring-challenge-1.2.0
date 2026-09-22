package output

type ProductSummary struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category Category  `json:"category"`
	Variants []Variant `json:"variants,omitempty"`
}
