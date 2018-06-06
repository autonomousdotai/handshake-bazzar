package response_obj

import (
	"time"

	"github.com/ninjadotorg/handshake-bazzar/bean"
	"github.com/ninjadotorg/handshake-bazzar/models"
	"github.com/ninjadotorg/handshake-bazzar/utils"
)

type ProductResponse struct {
	DateCreated   time.Time              `json:"date_created"`
	DateModified  time.Time              `json:"date_modified"`
	ID            int64                  `json:"id"`
	UserId        int64                  `json:"user_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Specification string                 `json:"specification"`
	Price         float64                `json:"price"`
	Shipping      int                    `json:"shipping"`
	Status        int                    `json:"status"`
	Images        []ProductImageResponse `json:"images"`
}

type ProductImageResponse struct {
	ID    int64  `json:"id"`
	Image string `json:"image"`
}

type ProductShakeResponse struct {
	ID       int64           `json:"id"`
	UserId   int64           `json:"user_id"`
	Price    float64         `json:"price"`
	Quantity int             `json:"quantity"`
	Amount   float64         `json:"amount"`
	Product  ProductResponse `json:"product"`
}

func MakeProductResponse(model models.Product) ProductResponse {
	result := ProductResponse{}
	result.ID = model.ID
	result.UserId = model.UserId
	result.Name = model.Name
	result.Description = model.Description
	result.Specification = model.Specification
	result.Price = model.Price
	result.Shipping = model.Shipping
	result.Status = model.Status
	result.Images = MakeArrayProductImageResponse(model.ProductImages)
	return result
}

func MakeProductImageResponse(model models.ProductImage) ProductImageResponse {
	result := ProductImageResponse{}
	result.ID = model.ID
	result.Image = utils.CdnUrlFor(model.Image)
	return result
}

func MakeArrayProductImageResponse(models []models.ProductImage) []ProductImageResponse {
	results := []ProductImageResponse{}
	for _, model := range models {
		result := MakeProductImageResponse(model)
		results = append(results, result)
	}
	return results
}

type UserResponse struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Status int    `json:"status"`
}

func MakeUserResponse(model models.User) UserResponse {
	result := UserResponse{}
	result.ID = model.ID
	result.Email = model.Email
	result.Name = model.Name
	result.Avatar = utils.CdnUrlFor(model.Avatar)
	return result
}

type ProductFaqResponse struct {
	ID           int64        `json:"id"`
	DateCreated  time.Time    `json:"date_created"`
	DateModified time.Time    `json:"date_modified"`
	UserId       int64        `json:"user_id"`
	ProductId    int64        `json:"product_id"`
	Question     string       `json:"question"`
	Answer       string       `json:"answer"`
	Status       int          `json:"status"`
	User         UserResponse `json:"user"`
}

func MakeProductFaqResponse(model models.ProductFaq) ProductFaqResponse {
	result := ProductFaqResponse{}
	result.ID = model.ID
	result.DateCreated = model.DateCreated
	result.DateModified = model.DateModified
	result.UserId = model.UserId
	result.ProductId = model.ProductId
	result.Question = model.Question
	result.Answer = model.Answer
	result.Status = model.Status
	result.User = MakeUserResponse(model.User)
	return result
}

func MakeArrayProductFaqResponse(models []models.ProductFaq) []ProductFaqResponse {
	results := []ProductFaqResponse{}
	for _, model := range models {
		result := MakeProductFaqResponse(model)
		results = append(results, result)
	}
	return results
}

func MakePaginationProductFaqResponse(pagination *bean.Pagination) PaginationResponse {
	return MakePaginationResponse(pagination.Page, pagination.PageSize, pagination.Total, MakeArrayProductFaqResponse(pagination.Items.([]models.ProductFaq)))
}
