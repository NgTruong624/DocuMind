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
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalysisResponse struct {
	FileHash       string   `json:"file_hash"`
	Summary        string   `json:"summary"`
	KeyClauses     []string `json:"key_clauses"`
	PotentialRisks []string `json:"potential_risks"`
}

type AnalysisListItem struct {
	ID             uint      `json:"id"`
	FileHash       string    `json:"file_hash"`
	CreatedAt      time.Time `json:"created_at"`
	SummaryPreview string    `json:"summary_preview"`
}

type AnalysisDetailResponse struct {
	ID             uint     `json:"id"`
	AnalysisID     uint     `json:"analysis_id"`
	Summary        string   `json:"summary"`
	KeyClauses     []string `json:"key_clauses"`
	PotentialRisks []string `json:"potential_risks"`
}

type ContractChatRequest struct {
	FileHash     string `json:"file_hash"`
	ContractText string `json:"contract_text"`
	Question     string `json:"question"`
}

type ContractChatResponse struct {
	Answer string `json:"answer"`
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
				FileHash:       existingAnalysis.FileHash,
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
	case strings.Contains(contentType, "officedocument.wordprocessingml.document") || strings.HasSuffix(filename, ".docx"):
		textContent, err = services.ExtractTextFromDOCX(fileReader)
	case strings.Contains(contentType, "msword") || strings.HasSuffix(filename, ".doc"):
		textContent, err = services.ExtractTextFromDOCX(fileReader)
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
		// Kiểm tra lỗi quota API
		if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "exceeded") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "API quota đã hết. Vui lòng thử lại sau hoặc liên hệ admin để nâng cấp quota."})
			return
		}
		
		// Kiểm tra lỗi API key
		if strings.Contains(err.Error(), "API_KEY") || strings.Contains(err.Error(), "authentication") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi xác thực API. Vui lòng kiểm tra cấu hình."})
			return
		}
		
		// Lỗi chung
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

	// 1. Tạo bản ghi chính (analyses) trước
	summaryPreview := analysisResp.Summary
	if len(summaryPreview) > 200 {
		summaryPreview = summaryPreview[:200]
	}
	analysisModel := models.Analysis{
		FileHash:       fileHash,
		SummaryPreview: summaryPreview,
	}
	if err := tx.Create(&analysisModel).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to save main analysis record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save main analysis record."})
		return
	}

	// 2. Tạo bản ghi chi tiết (analysis_details) với AnalysisID vừa tạo
	detail := models.AnalysisDetail{
		AnalysisID:     analysisModel.ID,
		Summary:        analysisResp.Summary,
		KeyClauses:     analysisResp.KeyClauses,
		PotentialRisks: analysisResp.PotentialRisks,
	}
	if err := tx.Create(&detail).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to save analysis details: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save analysis details."})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction."})
		return
	}

	c.JSON(http.StatusOK, AnalysisResponse{
		FileHash:       fileHash,
		Summary:        analysisResp.Summary,
		KeyClauses:     analysisResp.KeyClauses,
		PotentialRisks: analysisResp.PotentialRisks,
	})
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

// GET /api/v1/analyses - Lấy danh sách analyses (lịch sử)
func GetAnalyses(c *gin.Context) {
	var analyses []models.Analysis
	if err := database.DB.Order("created_at desc").Find(&analyses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analyses: " + err.Error()})
		return
	}
	var result []AnalysisListItem
	for _, a := range analyses {
		result = append(result, AnalysisListItem{
			ID:             a.ID,
			FileHash:       a.FileHash,
			CreatedAt:      a.CreatedAt,
			SummaryPreview: a.SummaryPreview,
		})
	}
	c.JSON(http.StatusOK, result)
}

// GET /api/v1/analyses/:id - Lấy chi tiết analysis
func GetAnalysisDetail(c *gin.Context) {
	id := c.Param("id")
	var detail models.AnalysisDetail
	if err := database.DB.Where("analysis_id = ?", id).First(&detail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Analysis detail not found"})
		return
	}
	resp := AnalysisDetailResponse{
		ID:             detail.ID,
		AnalysisID:     detail.AnalysisID,
		Summary:        detail.Summary,
		KeyClauses:     detail.KeyClauses,
		PotentialRisks: detail.PotentialRisks,
	}
	c.JSON(http.StatusOK, resp)
}

func ContractChatHandler(c *gin.Context) {
	var req ContractChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if utf8.RuneCountInString(req.Question) == 0 || utf8.RuneCountInString(req.Question) > 300 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Câu hỏi phải từ 1 đến 300 ký tự."})
		return
	}

	var contractText string
	if req.FileHash != "" {
		// Lấy nội dung hợp đồng từ DB
		var analysis models.Analysis
		result := database.DB.Where("file_hash = ?", req.FileHash).Preload("AnalysisDetail").First(&analysis)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy hợp đồng với file_hash đã cung cấp."})
			return
		}
		contractText = analysis.AnalysisDetail.Summary + "\n" + strings.Join(analysis.AnalysisDetail.KeyClauses, "\n") + "\n" + strings.Join(analysis.AnalysisDetail.PotentialRisks, "\n")
	} else if req.ContractText != "" {
		contractText = req.ContractText
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần cung cấp file_hash hoặc contract_text."})
		return
	}

	aiAnswer, err := services.AskContractQuestion(contractText, req.Question)
	if err != nil {
		// Kiểm tra lỗi quota API
		if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "exceeded") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "API quota đã hết. Vui lòng thử lại sau hoặc liên hệ admin để nâng cấp quota."})
			return
		}
		
		// Kiểm tra lỗi API key
		if strings.Contains(err.Error(), "API_KEY") || strings.Contains(err.Error(), "authentication") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi xác thực API. Vui lòng kiểm tra cấu hình."})
			return
		}
		
		// Lỗi chung
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI trả lời thất bại: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, ContractChatResponse{Answer: aiAnswer})
}