package itemsvc

type GetManyFilters struct {
	Seller   string
	Query    string
	Category string
}

type CreateUpdateInput struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	Price        float64  `json:"price"`
	Condition    int16    `json:"condition"`
	IsNegotiable bool     `json:"is_negotiable"`
	Category     string   `json:"category"`
	Subcategory  string   `json:"subcategory"`
}
