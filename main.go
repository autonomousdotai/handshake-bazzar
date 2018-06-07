package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-bazzar/api"
	"github.com/ninjadotorg/handshake-bazzar/configs"
	"google.golang.org/api/option"
)

func main() {

	configs.Initialize(os.Getenv("APP_CONF"))

	err := NewProcesser()
	if err != nil {
		panic(err)
	}

	// Logger
	logFile, err := os.OpenFile("logs/autonomous_service.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter) // You may need this
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// end Logger
	// Setting router
	router := gin.Default()
	router.Use(Logger())
	router.Use(AuthorizeMiddleware())
	// Router Index
	index := router.Group("/")
	{
		index.GET("/", func(context *gin.Context) {
			result := map[string]interface{}{
				"status":  1,
				"message": "Bazzar Service API",
			}
			context.JSON(http.StatusOK, result)
		})
	}
	productApi := api.ProductApi{}
	productApi.Init(router)

	faqApi := api.FaqApi{}
	faqApi.Init(router)

	router.Run(fmt.Sprintf(":%d", configs.AppConf.ServicePort))
}

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		context.Next()
		status := context.Writer.Status()
		latency := time.Since(t)
		log.Print("Request: " + context.Request.URL.String() + " | " + context.Request.Method + " - Status: " + strconv.Itoa(status) + " - " +
			latency.String())
	}
}

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, _ := strconv.ParseInt(context.GetHeader("User-Id"), 10, 64)
		context.Set("UserId", userId)
		context.Next()
	}
}

func NewProcesser() error {
	log.Println("NewProcesser")

	opt := option.WithCredentialsFile(configs.AppConf.PubsubConf.CredsFile)
	pubsubClient, err := pubsub.NewClient(context.Background(), configs.AppConf.PubsubConf.ProjectId, opt)
	if err != nil {
		log.Println(err)
		return err
	}

	handler, err := api.NewEthHandler(pubsubClient, configs.AppConf.PubsubConf.Topic, configs.AppConf.PubsubConf.Subscription)
	if err != nil {
		return err
	}

	go handler.Receive()
	return nil
}
