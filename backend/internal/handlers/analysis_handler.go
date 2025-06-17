package handlers

import (
	"documind/backend/internal/services" // Import service AI bạn vừa tạo
	"encoding/json"                      // Import thư viện để xử lý JSON
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnalysisResponse định nghĩa cấu trúc JSON mà chúng ta sẽ trả về cho frontend.
// Struct này phải khớp với cấu trúc JSON mà bạn yêu cầu AI trả về.
type AnalysisResponse struct {
	Summary        string   `json:"summary"`
	KeyClauses     []string `json:"key_clauses"`
	PotentialRisks []string `json:"potential_risks"`
}

// AnalyzeHandler xử lý yêu cầu phân tích file
func AnalyzeHandler(c *gin.Context) {
	// 1. Lấy và đọc file (Phần này giữ nguyên như cũ)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed: " + err.Error()})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file: " + err.Error()})
		return
	}
	defer file.Close()

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file content: " + err.Error()})
		return
	}
	textContent := string(contentBytes)

	// 2. Gọi AI Service để phân tích văn bản (Đây là logic mới)
	log.Println("Bắt đầu gọi AI để phân tích...")
	aiResultString, err := services.AnalyzeText(textContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed: " + err.Error()})
		return
	}

	// 3. Parse chuỗi JSON từ AI vào struct (Đây là logic mới)
	var analysisResp AnalysisResponse
	err = json.Unmarshal([]byte(aiResultString), &analysisResp)
	if err != nil {
		// Ghi log lỗi chi tiết hơn, bao gồm cả chuỗi AI trả về để dễ gỡ lỗi
		log.Printf("Lỗi khi parse JSON từ AI: %v. \nChuỗi nhận được: %s", err, aiResultString)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response."})
		return
	}

	// 4. Trả về kết quả có cấu trúc hoàn chỉnh cho frontend (Cập nhật)
	c.JSON(http.StatusOK, analysisResp)
}