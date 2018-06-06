package dao

import (
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
	"strings"
)

type EthTxDao struct {
}

func (ethTxDao EthTxDao) GetById(id int64) (models.EthTx) {
	dto := models.EthTx{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (ethTxDao EthTxDao) Create(dto models.EthTx, tx *gorm.DB) (models.EthTx, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.Hash = strings.ToLower(dto.Hash)
	dto.FromAddress = strings.ToLower(dto.FromAddress)
	dto.ToAddress = strings.ToLower(dto.ToAddress)
	dto.DateCreated = time.Now()
	dto.DateModified = dto.DateCreated
	err := tx.Create(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (ethTxDao EthTxDao) Update(dto models.EthTx, tx *gorm.DB) (models.EthTx, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.Hash = strings.ToLower(dto.Hash)
	dto.FromAddress = strings.ToLower(dto.FromAddress)
	dto.ToAddress = strings.ToLower(dto.ToAddress)
	dto.DateModified = time.Now()
	err := tx.Save(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (ethTxDao EthTxDao) Delete(dto models.EthTx, tx *gorm.DB) (models.EthTx, error) {
	if tx == nil {
		tx = models.Database()
	}
	err := tx.Delete(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}
