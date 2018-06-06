package bean

type Pagination struct {
	PageSize      int
	Page          int
	Total         int
	Items         interface{}
	Random        bool
	SearchRequest SearchRequest
}

type SearchRequest struct {
	QueryStr    string
	Country     string
	WarehouseId int
	Status      int
	OrderCode   string
}
