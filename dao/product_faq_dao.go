package dao

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-bazzar/bean"
	"github.com/ninjadotorg/handshake-bazzar/models"
)

type ProductFaqDao struct {
}

func (productFaqDao ProductFaqDao) GetById(id int64) models.ProductFaq {
	dto := models.ProductFaq{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (productFaqDao ProductFaqDao) Create(dto models.ProductFaq, tx *gorm.DB) (models.ProductFaq, error) {
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

func (productFaqDao ProductFaqDao) Update(dto models.ProductFaq, tx *gorm.DB) (models.ProductFaq, error) {
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

func (productFaqDao ProductFaqDao) Delete(dto models.ProductFaq, tx *gorm.DB) (models.ProductFaq, error) {
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

func (productFaqDao ProductFaqDao) GetAllBy(userId int64, productId int64, pagination *bean.Pagination) (*bean.Pagination, error) {
	dtos := []models.ProductFaq{}
	db := models.Database()
	if pagination != nil {
		db = db.Limit(pagination.PageSize)
		db = db.Offset(pagination.PageSize * (pagination.Page - 1))
	}
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}
	if productId > 0 {
		db = db.Where("product_id = ?", productId)
	}
	err := db.Order("prioriry asc, date_created desc").Find(&dtos).Error
	if err != nil {
		log.Print(err)
		return pagination, err
	}
	pagination.Items = dtos
	total := 0
	if pagination.Page == 1 && len(dtos) < pagination.PageSize {
		total = len(dtos)
	} else {
		err := db.Find(&dtos).Count(&total).Error
		if err != nil {
			log.Print(err)
			return pagination, err
		}
	}
	pagination.Total = total
	return pagination, nil
}
