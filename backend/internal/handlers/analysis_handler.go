package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler xử lý yêu cầu phân tích file
func AnalyzeHandler(c *gin.Context) {
	// 'file' là key mà frontend đã dùng trong formData.append('file', file)
	file, err := c.FormFile("file")

	// Xử lý lỗi nếu không có file hoặc tên field không đúng
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed: " + err.Error()})
		return
	}

	// Ghi log để xác nhận đã nhận được file
	log.Printf("Received file: %s", file.Filename)
	log.Printf("File size: %d bytes", file.Size)

	// Tạm thời trả về thông báo thành công
	// Sau này chúng ta sẽ thay thế phần này bằng logic gọi AI
	c.JSON(http.StatusOK, gin.H{
		"status":  "File received successfully",
		"fileName": file.Filename,
	})
}