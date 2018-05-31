package response_obj

import (
	"time"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/models"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/utils"
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
