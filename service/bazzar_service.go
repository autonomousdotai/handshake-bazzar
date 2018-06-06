package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-bazzar/configs"
	"github.com/ninjadotorg/handshake-bazzar/models"
	"github.com/ninjadotorg/handshake-bazzar/request_obj"
	"github.com/ninjadotorg/handshake-bazzar/utils"
	solr "github.com/rtt/Go-Solr"
)

type BazzarService struct {
}

func (crowdService BazzarService) CreateTx(userId int64, address string, hash string, refType string, refId int64, tx *gorm.DB) (models.EthTx, error) {
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
		return ethTx, err
	}
	return ethTx, nil
}

func (crowdService BazzarService) CreateProduct(userId int64, request request_obj.ProductRequest, context *gin.Context) (models.Product, error) {
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
		return product, err
	}

	imageLength, err := strconv.Atoi(context.Request.PostFormValue("image_length"))
	for i := 0; i < imageLength; i++ {
		imageFile, imageFileHeader, err := context.Request.FormFile(fmt.Sprintf("image_%d", i))
		if err != nil {
			log.Println(err)
			//rollback
			tx.Rollback()
			return product, err
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
				return product, err
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
			return product, err
		}
	}

	tx.Commit()

	product = productDao.GetFullById(product.ID)

	return product, nil
}

func (crowdService BazzarService) UpdateProduct(userId int64, crowdFundingId int64, request request_obj.ProductRequest, imageFile *multipart.File, imageFileHeader *multipart.FileHeader) (models.Product, error) {
	crowdFunding := productDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 || crowdFunding.UserId != userId {
		return crowdFunding, errors.New("crowdFundingId is invalid")
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
		return crowdFunding, err
	}
	return crowdFunding, nil
}

func (crowdService BazzarService) GetProduct(userId int64, crowdFundingId int64) (models.Product, error) {
	crowdFunding := productDao.GetById(crowdFundingId)
	if crowdFunding.ID <= 0 {
		return crowdFunding, errors.New("crowdFundingId is invalid")
	}
	return crowdFunding, nil
}

func (crowdService BazzarService) ShakeProduct(userId int64, productId int64, quantity int, address string, hash string) (models.ProductShake, error) {
	productShake := models.ProductShake{}

	if quantity <= 0 {
		return productShake, errors.New("quantity is invalid")
	}

	crowdFunding := productDao.GetFullById(productId)
	if crowdFunding.ID <= 0 {
		return productShake, errors.New("productId is invalid")
	}

	productShake.UserId = userId
	productShake.ProductId = productId
	productShake.Quantity = quantity
	productShake.Amount = float64(productShake.Quantity) * crowdFunding.Price
	productShake.Status = utils.ORDER_STATUS_SHAKED_PROCESS

	productShake, err := productShakeDao.Create(productShake, nil)
	if err != nil {
		log.Println(err)
		return productShake, err
	}

	_, err = crowdService.CreateTx(userId, address, hash, "payable_shake", productShake.ID, nil)
	if err != nil {
		log.Println(err)
		return productShake, err
	}

	return productShake, nil
}

func (crowdService BazzarService) DeliverProductShake(userId int64, productShakeId int64, address string, hash string) error {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return errors.New("crowdFunding is not shaked")
	}
	tx := models.Database().Begin()
	_, err := crowdService.CreateTx(userId, address, hash, "payable_deliver", userId, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}
	productShake.Status = utils.ORDER_STATUS_DELIVERED_PROCESS
	productShake, err = productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) CancelProductShake(userId int64, productShakeId int64, address string, hash string) error {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return errors.New("crowdFunding is not shaked")
	}
	tx := models.Database().Begin()
	_, err := crowdService.CreateTx(userId, address, hash, "payable_cancel", userId, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	productShake.Status = utils.ORDER_STATUS_CANCELED_PROCESS
	productShake, err = productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) RejectProductShake(userId int64, productShakeId int64, address string, hash string) error {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return errors.New("crowdFunding is not shaked")
	}
	tx := models.Database().Begin()
	_, err := crowdService.CreateTx(userId, address, hash, "payable_reject", userId, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}
	productShake.Status = utils.ORDER_STATUS_REJECTED_PROCESS
	productShake, err = productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) AcceptProductShake(userId int64, productShakeId int64, address string, hash string) error {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return errors.New("crowdFunding is not shaked")
	}
	tx := models.Database().Begin()
	_, err := crowdService.CreateTx(userId, address, hash, "payable_accept", userId, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}
	productShake.Status = utils.ORDER_STATUS_ACCEPTED_PROCESS
	productShake, err = productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) WithdrawProductShake(userId int64, productShakeId int64, address string, hash string) error {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 || productShake.Status <= 0 {
		return errors.New("crowdFunding is not shaked")
	}
	tx := models.Database().Begin()
	_, err := crowdService.CreateTx(userId, address, hash, "payable_withdraw", userId, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}
	productShake.Status = utils.ORDER_STATUS_WITHDRAWED_PROCESS
	productShake, err = productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)

		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) IndexSolr(productId int64) error {
	product := productDao.GetFullById(productId)

	crowdFundingImages := productImageDao.GetByProductId(product.ID)
	imageUrls := []string{}
	for _, crowdFundingImage := range crowdFundingImages {
		imageUrls = append(imageUrls, crowdFundingImage.Image)
	}

	document := map[string]interface{}{
		"add": []interface{}{
			map[string]interface{}{
				"id":                fmt.Sprintf("bazzar_%d", product.ID),
				"hid_l":             0,
				"type_i":            1,
				"state_i":           0,
				"init_user_id_i":    product.UserId,
				"shake_user_ids_is": []int64{},
				"text_search_ss":    []string{product.Name, product.Description, product.Specification},
				"shake_count_i":     0,
				"view_count_i":      0,
				"comment_count_i":   0,
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
	url := fmt.Sprintf("%s/%s", configs.AppConf.SolrServiceUrl, "handshake/update")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	bodyBytes, err := netUtil.CurlRequest(req)
	if err != nil {
		return err
	}
	result := solr.UpdateResponse{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (crowdService BazzarService) ProcessEventInit(hid int64, productShakeId int64) error {
	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		return errors.New("productShake is invalid")
	}

	productShake.Hid = hid

	productShake, err := productShakeDao.Update(productShake, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (crowdService BazzarService) ProcessEventShake(hid int64, productShakeId int64, fromAddress string) error {
	tx := models.Database().Begin()

	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		tx.Rollback()
		return errors.New("productShake is invalid")
	}

	productShake.Status = utils.ORDER_STATUS_SHAKED

	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	product := productDao.GetByHid(hid)
	if product.ID <= 0 {
		tx.Rollback()
		return errors.New("product is invalid")
	}

	product.ShakeNum += 1

	product, err = productDao.Update(product, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) ProcessEventDeliver(hid int64, productShakeId int64) error {
	tx := models.Database().Begin()

	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		tx.Rollback()
		return errors.New("productShake is invalid")
	}

	productShake.Status = utils.ORDER_STATUS_DELIVERED

	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) ProcessEventReject(hid int64, productShakeId int64) error {
	tx := models.Database().Begin()

	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		tx.Rollback()
		return errors.New("productShake is invalid")
	}

	productShake.Status = utils.ORDER_STATUS_REJECTED

	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) ProcessEventAccept(hid int64, productShakeId int64) error {
	tx := models.Database().Begin()

	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		tx.Rollback()
		return errors.New("productShake is invalid")
	}

	productShake.Status = utils.ORDER_STATUS_ACCEPTED

	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) ProcessEventCancel(hid int64, productShakeId int64) error {
	tx := models.Database().Begin()

	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		tx.Rollback()
		return errors.New("productShake is invalid")
	}

	productShake.Status = utils.ORDER_STATUS_CANCELED

	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crowdService BazzarService) ProcessEventWithdraw(hid int64, productShakeId int64) error {
	tx := models.Database().Begin()

	productShake := productShakeDao.GetById(productShakeId)
	if productShake.ID <= 0 {
		tx.Rollback()
		return errors.New("productShake is invalid")
	}

	productShake.Status = utils.ORDER_STATUS_WITHDRAWED

	productShake, err := productShakeDao.Update(productShake, tx)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
