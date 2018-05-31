package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type ProductImage struct {
	DateCreated  time.Time
	DateModified time.Time
	ID           int64
	ProductId    int64
	Image        string
	YoutubeUrl   string
	Priority     int
}

func (ProductImage) TableName() string {
	return "product_image"
}
