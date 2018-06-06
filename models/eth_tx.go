package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type EthTx struct {
	ID           int `gorm:"primary_key"`
	DateCreated  time.Time
	DateModified time.Time
	UserId       int64
	Hash         string
	RefType      string
	RefId        int64
	Status       int
	Value        float64
	FromAddress  string
	ToAddress    string
}

func (EthTx) TableName() string {
	return "eth_tx"
}
