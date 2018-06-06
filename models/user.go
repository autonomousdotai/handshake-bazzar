package models

import (
	_ "encoding/gob"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type JsonUserResponse struct {
	Status int  `json:"status"`
	Data   User `json:"data"`
}

type User struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Status int    `json:"status"`
}
