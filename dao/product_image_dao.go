package dao

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-bazzar/models"
)

type ProductImageDao struct {
}

func (productImageDao ProductImageDao) GetById(id int) models.ProductImage {
	dto := models.ProductImage{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (productImageDao ProductImageDao) Create(dto models.ProductImage, tx *gorm.DB) (models.ProductImage, error) {
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

func (productImageDao ProductImageDao) Update(dto models.ProductImage, tx *gorm.DB) (models.ProductImage, error) {
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

func (productImageDao ProductImageDao) Delete(dto models.ProductImage, tx *gorm.DB) (models.ProductImage, error) {
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

func (productImageDao ProductImageDao) GetByProductId(productId int64) []models.ProductImage {
	dtos := []models.ProductImage{}
	err := models.Database().Where("product_id = ?", productId).Find(&dtos).Error
	if err != nil {
		log.Print(err)
	}
	return dtos
}
