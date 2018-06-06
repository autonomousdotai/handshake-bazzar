package dao

import (
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type ProductShakeDao struct {
}

func (productShakeDao ProductShakeDao) GetById(id int64) (models.ProductShake) {
	dto := models.ProductShake{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (productShakeDao ProductShakeDao) Create(dto models.ProductShake, tx *gorm.DB) (models.ProductShake, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateCreated = time.Now()
	dto.DateModified = dto.DateCreated
	err := tx.Create(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (productShakeDao ProductShakeDao) Update(dto models.ProductShake, tx *gorm.DB) (models.ProductShake, error) {
	if tx == nil {
		tx = models.Database()
	}
	dto.DateModified = time.Now()
	err := tx.Save(&dto).Error
	if err != nil {
		log.Println(err)
		return dto, err
	}
	return dto, nil
}

func (productShakeDao ProductShakeDao) Delete(dto models.ProductShake, tx *gorm.DB) (models.ProductShake, error) {
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
