package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/farzadamr/event-manager-api/api/dto"
	"github.com/farzadamr/event-manager-api/config"
	"github.com/farzadamr/event-manager-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	config *config.Config
	logger logging.Logger
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	return &FileHandler{
		config: cfg,
		logger: logging.NewLogger(cfg),
	}
}

func (h *FileHandler) UploadPoster(c *gin.Context) {
	file, err := c.FormFile("poster")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "فایل یافت نشد"})
		return
	}

	if filepath.Ext(file.Filename) != ".jpg" && filepath.Ext(file.Filename) != ".jpeg" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "فقط فرمت JPG مجاز است"})
		return
	}
	fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
	folderPath := "storage/poster"

	_ = os.MkdirAll(folderPath, os.ModePerm)

	fullPath := filepath.Join(folderPath, fileName)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "خطا در ذخیره فایل"})
		return
	}

	c.JSON(http.StatusOK, dto.FileRef{
		Path: "/static/" + fileName,
		Mime: file.Header.Get("Content-Type"),
		Size: file.Size,
	})
}
