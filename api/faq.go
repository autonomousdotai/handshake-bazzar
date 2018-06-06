package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-bazzar/bean"
	"github.com/ninjadotorg/handshake-bazzar/request_obj"
	"github.com/ninjadotorg/handshake-bazzar/response_obj"
)

type FaqApi struct {
}

func (faqApi FaqApi) Init(router *gin.Engine) *gin.RouterGroup {
	faq := router.Group("/faq")
	{
		faq.GET("/:product_id", func(context *gin.Context) {
			context.String(200, "Common API")
		})
		faq.POST("/:product_id", func(context *gin.Context) {
			faqApi.CreateFaq(context)
		})
		faq.PUT("/:faq_id", func(context *gin.Context) {
			faqApi.CreateFaq(context)
		})
	}
	return faq
}

func (faqApi FaqApi) CreateFaq(context *gin.Context) {
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

	request := new(request_obj.ProductFaqRequest)
	err = context.Bind(&request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	faq, err := bazzarService.CreateFaq(userId.(int64), productId, *request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = response_obj.MakeProductFaqResponse(faq)
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}

func (faqApi FaqApi) UpdateFaq(context *gin.Context) {
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
	productFaqId, err := strconv.ParseInt(context.Param("faq_id"), 10, 64)
	if err != nil {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}
	if productFaqId <= 0 {
		result.SetStatus(bean.UnexpectedError)
		context.JSON(http.StatusOK, result)
		return
	}

	request := new(request_obj.ProductFaqRequest)
	err = context.Bind(&request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	faq, err := bazzarService.UpdateFaq(userId.(int64), productFaqId, *request)
	if err != nil {
		log.Print(err)
		result.SetStatus(bean.UnexpectedError)
		result.Error = err.Error()
		context.JSON(http.StatusOK, result)
		return
	}

	result.Data = response_obj.MakeProductFaqResponse(faq)
	result.Status = 1
	result.Message = ""
	context.JSON(http.StatusOK, result)
	return
}
