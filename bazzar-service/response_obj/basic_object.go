package response_obj

import "github.com/autonomousdotai/handshake-bazzar/bazzar-service/bean"

type ResponseObject struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

func (e *ResponseObject) SetStatus(key string) {
	e.Status = bean.CodeMessage[key].Code
	e.Message = bean.CodeMessage[key].Message
}

type PaginationResponse struct {
	PageSize int         `json:"page_size"`
	Page     int         `json:"page"`
	Total    int         `json:"total"`
	Items    interface{} `json:"items"`
	Status   int         `json:"status"`
	Message  string      `json:"message"`
}

func (e *PaginationResponse) SetStatus(key string) {
	e.Status = bean.CodeMessage[key].Code
	e.Message = bean.CodeMessage[key].Message
}

func MakePaginationResponse(page int, pageSize int, total int, items interface{}) PaginationResponse {
	result := PaginationResponse{}
	result.Page = page
	result.PageSize = pageSize
	result.Total = total
	result.Items = items
	return result
}
