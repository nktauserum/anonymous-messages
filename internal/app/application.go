package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nktauserum/anonymous-messages/config"
	"github.com/nktauserum/anonymous-messages/internal/bot"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
		AllowOrigins:     []string{"*"}, // или список ваших фронт-адресов
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	a.r.POST("/message", func(c *gin.Context) {
		tgbot, err := bot.LoadBot()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		message := c.PostForm("message")

		file, err := c.FormFile("file")
		if err != nil {
			if !errors.Is(err, http.ErrMissingFile) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			params := &telego.SendMessageParams{
				ChatID: tu.ID(a.c.Telegram.Admin),
				Text:   message,
			}

			_, err = tgbot.SendMessage(context.Background(), params)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			return
		}
		filename := filepath.Base(file.Filename)

		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer os.Remove(filename)

		img_file, err := os.Open(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		params := &telego.SendPhotoParams{
			ChatID:  tu.ID(a.c.Telegram.Admin),
			Photo:   telego.InputFile{File: img_file},
			Caption: message,
		}

		_, err = tgbot.SendPhoto(context.Background(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

	})

	return a.r.Run(fmt.Sprintf(":%d", a.c.Port))
}
