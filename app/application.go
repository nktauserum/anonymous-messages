package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	"github.com/nktauserum/anonymous-messages/article"
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

	bot, err := telego.NewBot(a.c.Telegram.Token, telego.WithDefaultLogger(false, true))
	if err != nil {
		return err
	}

	// TODO: вынести значение пути в конфиг
	service, err := article.NewArticleStorage("sqlite.db")
	if err != nil {
		return err
	}

	handler := Handler{
		b:       bot,
		service: service,
	}

	//Healthcheck
	a.r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})

	a.r.POST("/articles/add", handler.CreateArticle)
	a.r.GET("/articles/:uuid", handler.ReadArticle)
	a.r.PUT("/articles/update/:uuid", handler.UpdateArticle)
	a.r.GET("/articles/list", handler.ListArticles)
	a.r.DELETE("/articles/delete/:uuid", handler.DeleteArticle)

	a.r.POST("/message", handler.Message)

	return a.r.Run(fmt.Sprintf(":%d", a.c.Port))
}
