package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "time"
	"time"
)

type ProductFaq struct {
	ID           int
	DateCreated  time.Time
	DateModified time.Time
	ProductId    int
	Question     string
	Answer       string
	UserId       int
	Priority     int
	Status       int
}

func (ProductFaq) TableName() string {
	return "product_faq"
}
