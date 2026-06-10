package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qc-system/internal/service"
)

type OCRHandler struct {
	ocrService *service.OCRService
}

func NewOCRHandler() *OCRHandler {
	return &OCRHandler{
		ocrService: service.NewOCRService(),
	}
}

func (h *OCRHandler) Recognize(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		manualInput := c.PostForm("manual_input")
		if manualInput != "" {
			c.JSON(http.StatusOK, gin.H{
				"text":   manualInput,
				"source": "manual",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image file provided"})
		return
	}
	defer file.Close()

	imageData := make([]byte, 1024*1024)
	n, _ := file.Read(imageData)
	imageData = imageData[:n]

	manualInput := c.PostForm("manual_input")
	result, err := h.ocrService.RecognizeWithFallback(imageData, manualInput)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "OCR识别失败，请手动录入",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"text":    result,
		"source":  "ocr",
	})
}

func (h *OCRHandler) RecognizeMarkVerify(c *gin.Context) {
	var req struct {
		ImageData   []byte `json:"image_data"`
		ManualText string `json:"manual_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.ocrService.RecognizeWithFallback(req.ImageData, req.ManualText)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success":   false,
			"ocr_error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"text":    result,
	})
}
