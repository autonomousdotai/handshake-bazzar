package models

import (
	"time"
	_ "time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ProductFaq struct {
	ID           int64
	DateCreated  time.Time
	DateModified time.Time
	ProductId    int64
	Question     string
	Answer       string
	UserId       int64
	Priority     int
	Status       int
	User         User
}

func (ProductFaq) TableName() string {
	return "product_faq"
}
