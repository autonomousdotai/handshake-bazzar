package models

import (
	_ "encoding/gob"
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ProductShake struct {
	DateCreated     time.Time
	DateModified    time.Time
	ID              int64
	Hid             int64
	ChainId         int64
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
