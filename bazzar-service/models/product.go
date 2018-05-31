package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type Product struct {
	DateCreated   time.Time
	DateModified  time.Time
	ID            int64
	UserId        int64
	Name          string
	Description   string
	Specification string
	Price         float64
	Shipping      int
	Status        int
	ShakeNum      int
	CommentNum    int
	ProductImages []ProductImage
}

func (Product) TableName() string {
	return "product"
}
