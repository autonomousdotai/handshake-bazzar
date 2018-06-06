package dao

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-bazzar/models"
)

type ProductDao struct {
}

func (productDao ProductDao) GetById(id int64) models.Product {
	dto := models.Product{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (productDao ProductDao) GetFullById(id int64) models.Product {
	dto := models.Product{}
	db := models.Database()
	db = db.Preload("ProductImages")
	db = db.Where("id = ?", id)
	err := db.First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (productDao ProductDao) Create(dto models.Product, tx *gorm.DB) (models.Product, error) {
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

func (productDao ProductDao) Update(dto models.Product, tx *gorm.DB) (models.Product, error) {
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

func (productDao ProductDao) Delete(dto models.Product, tx *gorm.DB) (models.Product, error) {
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
