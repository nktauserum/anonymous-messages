package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/nktauserum/anonymous-messages/article"
	"github.com/nktauserum/anonymous-messages/config"
)

type Handler struct {
	b       *telego.Bot
	service article.ArticleService
}

func NewHandler(token string) (*Handler, error) {
	bot, err := telego.NewBot(token, telego.WithDefaultLogger(false, true))
	if err != nil {
		return nil, err
	}

	return &Handler{b: bot}, nil
}

func (h *Handler) Message(c *gin.Context) {
	conf := config.MustLoadConfig()

	message := fmt.Sprintf("%s\n\n<i>⌚️ %s</i>", c.PostForm("message"), time.Now().Format("15:04 02.01.2006"))

	file, err := c.FormFile("file")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Если изображения нет
		params := &telego.SendMessageParams{
			ChatID:    tu.ID(conf.Telegram.Admin),
			Text:      message,
			ParseMode: telego.ModeHTML,
		}

		_, err = h.b.SendMessage(context.Background(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		return
	}
	// Если изображение есть
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
	defer img_file.Close()

	params := &telego.SendPhotoParams{
		ChatID:                tu.ID(conf.Telegram.Admin),
		Photo:                 telego.InputFile{File: img_file},
		Caption:               message,
		ParseMode:             telego.ModeHTML,
		ShowCaptionAboveMedia: true,
	}

	_, err = h.b.SendPhoto(context.Background(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (h *Handler) CreateArticle(c *gin.Context) {
	var article article.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load time location"})
		return
	}

	article.DatePublished = time.Now().In(loc)
	article.UUID = uuid.New().String()

	if err := h.service.CreateArticle(context.Background(), article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Article created successfully"})
}

func (h *Handler) ListArticles(c *gin.Context) {
	articles, err := h.service.ListArticles(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *Handler) UpdateArticle(c *gin.Context) {
	articleUUID := c.Param("uuid")
	var request struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateArticle(context.Background(), articleUUID, request.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

func (h *Handler) ReadArticle(c *gin.Context) {
	articleUUID := c.Param("uuid")
	article, err := h.service.ReadArticle(context.Background(), articleUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read article"})
		return
	}
	article.UUID = articleUUID

	c.JSON(http.StatusOK, article)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	articleUUID := c.Param("uuid")
	if err := h.service.DeleteArticle(context.Background(), articleUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
