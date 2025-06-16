package handlers

import (
	"io" // Thêm thư viện io
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler xử lý yêu cầu phân tích file
func AnalyzeHandler(c *gin.Context) {
	// 1. Lấy file từ request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed: " + err.Error()})
		return
	}

	// 2. Mở file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file: " + err.Error()})
		return
	}
	defer file.Close() // Đảm bảo file được đóng sau khi hàm kết thúc

	// 3. Đọc toàn bộ nội dung của file
	contentBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file content: " + err.Error()})
		return
	}

	// Chuyển nội dung từ dạng bytes sang string
	textContent := string(contentBytes)

	// In nội dung ra log để kiểm tra
	log.Printf("Extracted text: %s", textContent)

	// 4. Trả nội dung về cho frontend (tạm thời để trong summary)
	c.JSON(http.StatusOK, gin.H{
		"summary":        textContent,
		"key_clauses":    []string{}, // Trả về mảng rỗng
		"potential_risks": []string{}, // Trả về mảng rỗng
	})
}