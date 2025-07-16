package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nktauserum/anonymous-messages/config"
)

type Application struct {
	r *gin.Engine
	c config.Config
}

func NewApplication() *Application {
	return &Application{c: *config.MustLoadConfig(), r: gin.Default()}
}

func (a *Application) Run() error {
	log.Println("started!")

	a.r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler, err := NewHandler(a.c.Telegram.Token)
	if err != nil {
		return err
	}

	//Healthcheck
	a.r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})

	a.r.POST("/message", handler.Message)

	return a.r.Run(fmt.Sprintf(":%d", a.c.Port))
}
