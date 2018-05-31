package dao

import (
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/models"
	"log"
	"github.com/jinzhu/gorm"
	"time"
)

type CrowdFundingFaqDao struct {
}

func (crowdFundingFaqDao CrowdFundingFaqDao) GetById(id int) (models.ProductFaq) {
	dto := models.ProductFaq{}
	err := models.Database().Where("id = ?", id).First(&dto).Error
	if err != nil {
		log.Print(err)
	}
	return dto
}

func (crowdFundingFaqDao CrowdFundingFaqDao) Create(dto models.ProductFaq, tx *gorm.DB) (models.ProductFaq, error) {
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

func (crowdFundingFaqDao CrowdFundingFaqDao) Update(dto models.ProductFaq, tx *gorm.DB) (models.ProductFaq, error) {
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

func (crowdFundingFaqDao CrowdFundingFaqDao) Delete(dto models.ProductFaq, tx *gorm.DB) (models.ProductFaq, error) {
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
