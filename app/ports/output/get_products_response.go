package output

type GetAllProductsResponse struct {
	Products []Product `json:"products"`
}

type GetByProductCodeResponse struct {
	Product Product `json:"product"`
}
