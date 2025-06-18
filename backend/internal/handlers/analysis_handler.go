package handlers

import (
	"documind/backend/internal/models"
	"documind/backend/internal/services"
	"documind/backend/pkg/database"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

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

	// Detect file type using Content-Type AND file extension as a fallback
	contentType := fileHeader.Header.Get("Content-Type")
	filename := strings.ToLower(fileHeader.Filename)
	var textContent string

	switch {
	// Kiểm tra PDF
	case strings.Contains(contentType, "pdf") || strings.HasSuffix(filename, ".pdf"):
		textContent, err = services.ExtractTextFromPDF(file)
		if err != nil {
			log.Printf("Error extracting text from PDF '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể trích xuất nội dung từ file PDF."})
			return
		}

	// Kiểm tra DOCX
	case strings.Contains(contentType, "officedocument.wordprocessingml.document") || strings.HasSuffix(filename, ".docx"):
		textContent, err = services.ExtractTextFromDOCX(file)
		if err != nil {
			log.Printf("Error extracting text from DOCX '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể trích xuất nội dung từ file DOCX."})
			return
		}

	// Kiểm tra DOC
	case strings.Contains(contentType, "msword") || strings.HasSuffix(filename, ".doc"):
		textContent, err = services.ExtractTextFromDOCX(file)
		if err != nil {
			log.Printf("Error extracting text from DOC '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể trích xuất nội dung từ file DOC."})
			return
		}

	// Kiểm tra TXT
	case strings.Contains(contentType, "text/plain") || strings.HasSuffix(filename, ".txt"):
		var contentBytes []byte
		contentBytes, err = io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading plain text file '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc nội dung từ file text."})
			return
		}
		textContent = string(contentBytes)

	// Trường hợp không hỗ trợ
	default:
		log.Printf("Unsupported file format. Filename: '%s', Content-Type: '%s'", filename, contentType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng file không được hỗ trợ. Vui lòng thử file PDF, DOCX, hoặc DOC."})
		return
	}

	// 2. Gọi AI Service để phân tích văn bản
	log.Println("Bắt đầu gọi AI để phân tích...")
	aiResultString, err := services.AnalyzeText(textContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed: " + err.Error()})
		return
	}

	// 3. Parse chuỗi JSON từ AI vào struct
	cleanedJSONString := cleanAIResponse(aiResultString)

	var analysisResp AnalysisResponse
	err = json.Unmarshal([]byte(cleanedJSONString), &analysisResp)
	if err != nil {
		log.Printf("Lỗi khi parse JSON từ AI: %v. \nChuỗi nhận được: %s", err, aiResultString)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response."})
		return
	}

	// 4. Persist the result to the database
	analysisModel := models.Analysis{
		Summary:        analysisResp.Summary,
		KeyClauses:     analysisResp.KeyClauses,
		PotentialRisks: analysisResp.PotentialRisks,
	}
	if err := database.DB.Create(&analysisModel).Error; err != nil {
		log.Printf("Failed to save analysis to DB: %v", err)
	}

	// 5. Trả về kết quả có cấu trúc hoàn chỉnh cho frontend
	c.JSON(http.StatusOK, analysisResp)
}

func cleanAIResponse(s string) string {
	// Loại bỏ ```json ở đầu và ``` ở cuối
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```json") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimSuffix(s, "```")
	}
	// Đôi khi AI chỉ trả về ``` mà không có chữ json
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
	}
	return strings.TrimSpace(s)
}