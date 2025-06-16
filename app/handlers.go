package app

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nktauserum/anonymous-messages/bot"
	"github.com/nktauserum/anonymous-messages/config"
	"net/http"
	"os"
	"path/filepath"
)

func Message(c *gin.Context) {
	conf := config.MustLoadConfig()

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
			ChatID: tu.ID(conf.Telegram.Admin),
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
		ChatID:  tu.ID(conf.Telegram.Admin),
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

}
