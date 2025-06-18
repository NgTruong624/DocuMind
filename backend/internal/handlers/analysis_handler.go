package handlers

import (
	"crypto/sha256"
	"documind/backend/internal/models"
	"documind/backend/internal/services"
	"documind/backend/pkg/database"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalysisResponse struct {
	Summary        string   `json:"summary"`
	KeyClauses     []string `json:"key_clauses"`
	PotentialRisks []string `json:"potential_risks"`
}

func AnalyzeHandler(c *gin.Context) {
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

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file content: " + err.Error()})
		return
	}

	hash := sha256.New()
	hash.Write(fileBytes)
	fileHash := hex.EncodeToString(hash.Sum(nil))

	var existingAnalysis models.Analysis
	// Dùng Preload để GORM tự động lấy dữ liệu từ bảng analysis_details liên quan
	result := database.DB.Where("file_hash = ?", fileHash).Preload("AnalysisDetail").First(&existingAnalysis)

	// Cache Hit: Nếu tìm thấy và không có lỗi nào khác ngoài "không tìm thấy"
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if result.Error == nil {
			log.Printf("Cache hit for file hash: %s", fileHash)
			c.JSON(http.StatusOK, AnalysisResponse{
				Summary:        existingAnalysis.AnalysisDetail.Summary,
				KeyClauses:     existingAnalysis.AnalysisDetail.KeyClauses,
				PotentialRisks: existingAnalysis.AnalysisDetail.PotentialRisks,
			})
			return
		}
		// Xử lý các lỗi database khác nếu có
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + result.Error.Error()})
		return
	}

	// Cache Miss: Tiếp tục xử lý file mới
	log.Printf("Cache miss for file hash: %s. Processing new file.", fileHash)

	fileReader := strings.NewReader(string(fileBytes))
	contentType := fileHeader.Header.Get("Content-Type")
	filename := strings.ToLower(fileHeader.Filename)
	var textContent string

	// ... (Toàn bộ khối switch để trích xuất text vẫn giữ nguyên)
	switch {
	case strings.Contains(contentType, "pdf") || strings.HasSuffix(filename, ".pdf"):
		textContent, err = services.ExtractTextFromPDF(fileReader)
	// ... (các case khác)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng file không được hỗ trợ."})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract text from file: " + err.Error()})
		return
	}

	aiResultString, err := services.AnalyzeText(textContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed: " + err.Error()})
		return
	}

	cleanedJSONString := cleanAIResponse(aiResultString)
	var analysisResp AnalysisResponse
	err = json.Unmarshal([]byte(cleanedJSONString), &analysisResp)
	if err != nil {
		log.Printf("Lỗi khi parse JSON từ AI: %v. \nChuỗi gốc: %s", err, aiResultString)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response."})
		return
	}

	// SỬ DỤNG DATABASE TRANSACTION ĐỂ LƯU DỮ LIỆU VÀO 2 BẢNG
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction."})
		return
	}

	// 1. Tạo bản ghi chi tiết trước
	detail := models.AnalysisDetail{
		Summary:        analysisResp.Summary,
		KeyClauses:     analysisResp.KeyClauses,
		PotentialRisks: analysisResp.PotentialRisks,
	}
	if err := tx.Create(&detail).Error; err != nil {
		tx.Rollback() // Nếu lỗi, hủy bỏ transaction
		log.Printf("Failed to save analysis details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save analysis details."})
		return
	}

	// 2. Tạo bản ghi chính, liên kết với bản ghi chi tiết
	analysisModel := models.Analysis{
		FileHash: fileHash,
		// GORM sẽ tự động điền AnalysisID vào bảng details dựa trên liên kết
	}
	if err := tx.Model(&analysisModel).Association("AnalysisDetail").Append(&detail); err != nil {
		tx.Rollback()
		log.Printf("Failed to create analysis association: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create analysis association."})
		return
	}
	if err := tx.Create(&analysisModel).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to save main analysis record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save main analysis record."})
		return
	}


	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction."})
		return
	}

	c.JSON(http.StatusOK, analysisResp)
}

func cleanAIResponse(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimSuffix(s, "```")
	}
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}