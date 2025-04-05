package dto

type ProductOption struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ProductId string `json:"product_id"`
	Price     string `json:"price"`
	Metadata  string `json:"metadata"`
}

type ProductOptionDetail struct{}
