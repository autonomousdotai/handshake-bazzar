package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-bazzar/bean"
	"github.com/ninjadotorg/handshake-bazzar/request_obj"
	"github.com/ninjadotorg/handshake-bazzar/response_obj"
)

type ProductApi struct {
}

func (api ProductApi) Init(router *gin.Engine) *gin.RouterGroup {
	bazzarGroup := router.Group("/product")
	{
		bazzarGroup.POST("/", func(context *gin.Context) {
			api.CreateProduct(context)
		})
		bazzarGroup.PUT("/", func(context *gin.Context) {
			api.UpdateProduct(context)
		})
		bazzarGroup.GET("/:product_id", func(context *gin.Context) {
			api.GetProduct(context)
		})
		bazzarGroup.POST("/shake/:product_id", func(context *gin.Context) {
			api.ShakeProduct(context)
		})
		bazzarGroup.POST("/deliver/:product_shake_id", func(context *gin.Context) {
			api.DeliverProductShake(context)
		})
		bazzarGroup.POST("/cancel/:product_shake_id", func(context *gin.Context) {
			api.CancelProductShake(context)
		})
		bazzarGroup.POST("/reject/:product_shake_id", func(context *gin.Context) {
			api.RejectProductShake(context)
		})
		bazzarGroup.POST("/accept/:product_shake_id", func(context *gin.Context) {
			api.AcceptProductShake(context)
		})
		bazzarGroup.POST("/withdraw/:product_shake_id", func(context *gin.Context) {
			api.WithdrawProductShake(context)
		})
	}
	return bazzarGroup
}

func (self ProductApi) CreateProduct(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}

	requestJson := context.Request.PostFormValue("request")
	request := new(request_obj.ProductRequest)
	err := json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	crowdFunging, err := bazzarService.CreateProduct(userId.(int64), *request, context)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeProductResponse(crowdFunging)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) UpdateProduct(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productId, err := strconv.ParseInt(context.Param("product_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	requestJson := context.Request.PostFormValue("request")
	request := new(request_obj.ProductRequest)
	err = json.Unmarshal([]byte(requestJson), &request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	imageFile, imageFileHeader, err := context.Request.FormFile("image")
	product, err := bazzarService.UpdateProduct(userId.(int64), productId, *request, &imageFile, imageFileHeader)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeProductResponse(product)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) GetProduct(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	crowdFungingId, err := strconv.ParseInt(context.Param("crowd_funding_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if crowdFungingId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	crowdFunging, err := bazzarService.GetProduct(0, crowdFungingId)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	data := response_obj.MakeProductResponse(crowdFunging)

	result.Data = data
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) ShakeProduct(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productId, err := strconv.ParseInt(context.Param("product_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	quantity, err := strconv.Atoi(context.Query("quantity"))
	address := context.Query("address")
	hash := context.Query("hash")

	productShake, err := bazzarService.ShakeProduct(userId.(int64), productId, quantity, address, hash)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}
	_ = productShake

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) DeliverProductShake(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productShakeId, err := strconv.ParseInt(context.Param("product_shake_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productShakeId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	address := context.Query("address")
	hash := context.Query("hash")

	err = bazzarService.DeliverProductShake(userId.(int64), productShakeId, address, hash)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) CancelProductShake(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productShakeId, err := strconv.ParseInt(context.Param("product_shake_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productShakeId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	address := context.Query("address")
	hash := context.Query("hash")

	err = bazzarService.CancelProductShake(userId.(int64), productShakeId, address, hash)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) RejectProductShake(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productShakeId, err := strconv.ParseInt(context.Param("product_shake_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productShakeId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	address := context.Query("address")
	hash := context.Query("hash")

	err = bazzarService.RejectProductShake(userId.(int64), productShakeId, address, hash)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) AcceptProductShake(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productShakeId, err := strconv.ParseInt(context.Param("product_shake_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productShakeId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	address := context.Query("address")
	hash := context.Query("hash")

	err = bazzarService.AcceptProductShake(userId.(int64), productShakeId, address, hash)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (self ProductApi) WithdrawProductShake(context *gin.Context) {
	result := new(response_obj.ResponseObject)

	userId, ok := context.Get("UserId")
	if !ok {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	if userId.(int64) <= 0 {
		result.SetStatus(bean.NotSignIn)
		context.JSON(http.StatusOK, result)
		return
	}
	productShakeId, err := strconv.ParseInt(context.Param("product_shake_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productShakeId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	address := context.Query("address")
	hash := context.Query("hash")

	err = bazzarService.WithdrawProductShake(userId.(int64), productShakeId, address, hash)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}
