package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type ProductShake struct {
	DateCreated     time.Time
	DateModified    time.Time
	ID              int64
	UserId          int64
	ProductId       int64
	Price           float64
	Quantity        int
	Amount          float64
	Status          int
	Email           string
	ShippingAddress string
	Address         string
}

func (ProductShake) TableName() string {
	return "product_shake"
}
