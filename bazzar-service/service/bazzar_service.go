package service

import (
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/models"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/bean"
	"errors"
	"mime/multipart"
	"strings"
	"time"
	"log"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/request_obj"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/utils"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/autonomousdotai/handshake-bazzar/bazzar-service/configs"
)

type BazzarService struct {
}

func (crowdService BazzarService) CreateTx(userId int64, address string, hash string, refType string, refId int64, tx *gorm.DB) (models.EthTx, *bean.AppError) {
	ethTx := models.EthTx{}
	ethTx.UserId = userId
	ethTx.FromAddress = address
	ethTx.Hash = hash
	ethTx.RefType = refType
	ethTx.RefId = refId
	ethTx.Status = 0
	ethTx, err := ethTxDao.Create(ethTx, tx)
	if err != nil {
		log.Println(err)
		return ethTx, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}
	return ethTx, nil
}

func (crowdService BazzarService) CreateProduct(userId int64, request request_obj.ProductRequest, context *gin.Context) (models.Product, *bean.AppError) {
	product := models.Product{}

	tx := models.Database().Begin()

	product.UserId = userId
	product.Name = request.Name
	product.Price = request.Price
	product.Shipping = request.Shipping
	product.Status = 1

	product, err := productDao.Create(product, tx)
	if err != nil {
		log.Println(err)
		//rollback
		tx.Rollback()
		return product, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	imageLength, err := strconv.Atoi(context.Request.PostFormValue("image_length"))
	for i := 0; i < imageLength; i++ {
		imageFile, imageFileHeader, err := context.Request.FormFile(fmt.Sprintf("image_%d", i))
		if err != nil {
			log.Println(err)
			//rollback
			tx.Rollback()
			return product, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
		filePath := ""
		if imageFile != nil && imageFileHeader != nil {
			uploadImageFolder := "product"
			fileName := imageFileHeader.Filename
			imageExt := strings.Split(fileName, ".")[1]
			fileNameImage := fmt.Sprintf("product-%d-image-%s.%s", product.ID, time.Now().Format("20060102150405"), imageExt)
			filePath = uploadImageFolder + "/" + fileNameImage
			err := fileUploadService.UploadFile(filePath, &imageFile)
			if err != nil {
				log.Println(err)
				//rollback
				tx.Rollback()
				return product, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
			}
		}
		productImage := models.ProductImage{}

		productImage.ProductId = product.ID
		productImage.Image = filePath

		productImage, err = productImageDao.Create(productImage, tx)
		if err != nil {
			log.Println(err)
			//rollback
			tx.Rollback()
			return product, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
		}
	}

	tx.Commit()

	product = productDao.GetFullById(product.ID)

	return product, nil
}

func (crowdService BazzarService) UpdateProduct(userId int64, crowdFundingId int64, request request_obj.ProductRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.Product, *bean.AppError) {
	crowdFunding := productDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 || crowdFunding.UserId != userId {
		return crowdFunding, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}

	crowdFunding.Name = request.Name
	crowdFunding.Description = request.Description
	crowdFunding.Specification = request.Specification

	if crowdFunding.Status == 1 {
		crowdFunding.Price = request.Price
		crowdFunding.Shipping = request.Shipping
	}

	crowdFunding, err := productDao.Update(crowdFunding, nil)
	if err != nil {
		log.Println(err)
		return crowdFunding, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}
	return crowdFunding, nil
}

func (crowdService BazzarService) GetProduct(userId int64, crowdFundingId int64) (models.Product, *bean.AppError) {
	crowdFunding := productDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFunding, &bean.AppError{errors.New("crowdFundingId is invalid"), "crowdFundingId is invalid", -1, "error_occurred"}
	}
	return crowdFunding, nil
}

func (crowdService BazzarService) ShakeProduct(userId int64, productId int64, quantity int, address string, hash string) (models.ProductShake, *bean.AppError) {
	productShake := models.ProductShake{}

	if quantity <= 0 {
		return productShake, &bean.AppError{errors.New("quantity is invalid"), "quantity is invalid", -1, "error_occurred"}
	}

	crowdFunding := productDao.GetFullById(productId)
	if crowdFunding.ID <= 0 {
		return productShake, &bean.AppError{errors.New("productId is invalid"), "productId is invalid", -1, "error_occurred"}
	}

	productShake.UserId = userId
	productShake.ProductId = productId
	productShake.Quantity = quantity
	productShake.Amount = float64(productShake.Quantity) * crowdFunding.Price
	productShake.Status = utils.ORDER_STATUS_SHAKED_PROCESS

	productShake, err := productShakeDao.Create(productShake, nil)
	if err != nil {
		log.Println(err)
		return productShake, &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	_, appErr := crowdService.CreateTx(userId, address, hash, "payable_shake", productShake.ID, nil)
	if appErr != nil {
		log.Println(appErr.OrgError)
		return productShake, appErr
	}

	return productShake, nil
}

func (crowdService BazzarService) DeliverProductShake(userId int64, productShakeId int64, address string, hash string) (*bean.AppError) {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}
	tx := models.Database().Begin()
	_, appErr := crowdService.CreateTx(userId, address, hash, "payable_deliver", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}
	productShake.Status = utils.ORDER_STATUS_DELIVERED_PROCESS
	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) CancelProductShake(userId int64, productShakeId int64, address string, hash string) (*bean.AppError) {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}
	tx := models.Database().Begin()
	_, appErr := crowdService.CreateTx(userId, address, hash, "payable_cancel", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}
	productShake.Status = utils.ORDER_STATUS_CANCELED_PROCESS
	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) RejectProductShake(userId int64, productShakeId int64, address string, hash string) (*bean.AppError) {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}
	tx := models.Database().Begin()
	_, appErr := crowdService.CreateTx(userId, address, hash, "payable_reject", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}
	productShake.Status = utils.ORDER_STATUS_REJECTED_PROCESS
	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) AcceptProductShake(userId int64, productShakeId int64, address string, hash string) (*bean.AppError) {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}
	tx := models.Database().Begin()
	_, appErr := crowdService.CreateTx(userId, address, hash, "payable_accept", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}
	productShake.Status = utils.ORDER_STATUS_ACCEPTED_PROCESS
	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) WithdrawProductShake(userId int64, productShakeId int64, address string, hash string) (*bean.AppError) {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return &bean.AppError{errors.New("crowdFunding is not shaked"), "crowdFunding is not shaked", -1, "error_occurred"}
	}
	tx := models.Database().Begin()
	_, appErr := crowdService.CreateTx(userId, address, hash, "payable_withdraw", userId, tx)
	if appErr != nil {
		log.Println(appErr.OrgError)

		tx.Rollback()
		return appErr
	}
	productShake.Status = utils.ORDER_STATUS_WITHDRAWED_PROCESS
	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return &bean.AppError{errors.New(err.Error()), "Error occurred, please try again", -1, "error_occurred"}
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) MakeObjectToIndex(productId int64) (error) {
	product := productDao.GetFullById(productId)

	crowdFundingImages := productImageDao.GetByProductId(product.ID)
	imageUrls := []string{}
	for _, crowdFundingImage := range crowdFundingImages {
		imageUrls = append(imageUrls, crowdFundingImage.Image)
	}

	document := map[string]interface{}{
		"add": [] interface{}{
			map[string]interface{}{
				"id":                fmt.Sprintf("bazzar_%d", product.ID),
				"hid_l":             0,
				"type_i":            1,
				"state_i":           0,
				"init_user_id_i":    product.UserId,
				"shake_user_ids_is": []int64{},
				"text_search_ss":    []string{product.Name, product.Description, product.Specification},
				"shake_count_i":     product.ShakeNum,
				"view_count_i":      0,
				"comment_count_i":   product.CommentNum,
				"is_private_i":      0,
				"init_at_i":         product.DateCreated.Unix(),
				"last_update_at_i":  product.DateModified.Unix(),
				//custom fileds
				"name_s":              product.Name,
				"short_description_s": product.Description,
				"shipping_i":          product.Shipping,
				"image_ss":            imageUrls,
			},
		},
	}

	jsonStr, err := json.Marshal(document)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", configs.SolrServiceUrl+"/handshake/update", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	bodyBytes, err := netUtil.CurlRequest(req)
	if err != nil {
		return err
	}
	result := algoliasearch.BatchRes{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}