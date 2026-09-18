package internal

type Response struct {
	Detail string `json:"detail"`
}

type Offer struct {
	Store string `json:"store"`
	ProductName string `json:"product_name"`
	ProductDesc string `json:"product_desc"`
	Quantity string `json:"quantity"`
	Price float64 `json:"price"`
	PricePKU float64 `json:"price_per_kilo_unit"`
}


