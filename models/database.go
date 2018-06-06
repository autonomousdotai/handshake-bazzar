package models

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-bazzar/configs"
)

var databaseConn *gorm.DB = nil

func Database() *gorm.DB {
	//open a db connection
	if databaseConn == nil {
		d, err := gorm.Open("mysql", configs.AppConf.DbUrl)
		d.LogMode(true)
		if err != nil {
			log.Println(err)
			return nil
		}
		// skip save associations of gorm -> manual save by code
		databaseConn = d.Set("gorm:save_associations", false)
		databaseConn.DB().SetMaxOpenConns(20)
		databaseConn.DB().SetMaxIdleConns(10)
	}
	return databaseConn
}
