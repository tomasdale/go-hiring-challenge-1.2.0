package output

type ListCategoriesResponse struct {
	Categories []Category `json:"categories"`
}
type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
