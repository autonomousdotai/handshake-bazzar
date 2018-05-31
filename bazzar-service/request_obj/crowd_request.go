package request_obj

type ProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Specification string  `json:"specification"`
	Price         float64 `json:"price"`
	Shipping      int     `json:"shipping"`
	Status        int     `json:"status"`
}

type ProductFaqRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
