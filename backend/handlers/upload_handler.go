package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadDir string
	publicURL string
}

func NewUploadHandler() *UploadHandler {
	dir := "./uploads"
	os.MkdirAll(dir, 0755)
	return &UploadHandler{
		uploadDir: dir,
		publicURL: "/uploads",
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未获取到上传文件"})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
	}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 jpg/jpeg/png/gif/webp/bmp 图片格式"})
		return
	}

	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过 10MB"})
		return
	}

	filename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), file.Size, ext)
	savePath := filepath.Join(h.uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败: " + err.Error()})
		return
	}

	url := fmt.Sprintf("%s/%s", h.publicURL, filename)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"url":      url,
			"filename": filename,
			"size":     file.Size,
		},
	})
}

func (h *UploadHandler) UploadDir() string {
	return h.uploadDir
}

func (h *UploadHandler) PublicURL() string {
	return h.publicURL
}
