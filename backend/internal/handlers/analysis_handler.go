package handlers

import (
	"crypto/sha256"
	"documind/backend/internal/models"
	"documind/backend/internal/services"
	"documind/backend/pkg/database"
	"encoding/hex"
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

	// Đọc toàn bộ nội dung file vào bộ nhớ
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file: " + err.Error()})
		return
	}

	// Tính toán hash của file
	hasher := sha256.New()
	hasher.Write(fileBytes)
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Kiểm tra xem file đã được phân tích trước đó chưa
	var existingAnalysis models.Analysis
	result := database.DB.Where("file_hash = ?", fileHash).First(&existingAnalysis)
	if result.Error == nil {
		// File đã được phân tích trước đó, trả về kết quả cũ
		analysisResp := AnalysisResponse{
			Summary:        existingAnalysis.Summary,
			KeyClauses:     existingAnalysis.KeyClauses,
			PotentialRisks: existingAnalysis.PotentialRisks,
		}
		c.JSON(http.StatusOK, analysisResp)
		return
	}

	// Nếu file chưa được phân tích, tiếp tục xử lý
	contentType := fileHeader.Header.Get("Content-Type")
	filename := strings.ToLower(fileHeader.Filename)
	var textContent string

	// Tạo một io.Reader từ fileBytes để sử dụng lại
	fileReader := strings.NewReader(string(fileBytes))

	switch {
	// Kiểm tra PDF
	case strings.Contains(contentType, "pdf") || strings.HasSuffix(filename, ".pdf"):
		textContent, err = services.ExtractTextFromPDF(fileReader)
		if err != nil {
			log.Printf("Error extracting text from PDF '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể trích xuất nội dung từ file PDF."})
			return
		}

	// Kiểm tra DOCX
	case strings.Contains(contentType, "officedocument.wordprocessingml.document") || strings.HasSuffix(filename, ".docx"):
		textContent, err = services.ExtractTextFromDOCX(fileReader)
		if err != nil {
			log.Printf("Error extracting text from DOCX '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể trích xuất nội dung từ file DOCX."})
			return
		}

	// Kiểm tra DOC
	case strings.Contains(contentType, "msword") || strings.HasSuffix(filename, ".doc"):
		textContent, err = services.ExtractTextFromDOCX(fileReader)
		if err != nil {
			log.Printf("Error extracting text from DOC '%s': %v", filename, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể trích xuất nội dung từ file DOC."})
			return
		}

	// Kiểm tra TXT
	case strings.Contains(contentType, "text/plain") || strings.HasSuffix(filename, ".txt"):
		textContent = string(fileBytes)

	// Trường hợp không hỗ trợ
	default:
		log.Printf("Unsupported file format. Filename: '%s', Content-Type: '%s'", filename, contentType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng file không được hỗ trợ. Vui lòng thử file PDF, DOCX, hoặc DOC."})
		return
	}

	// Gọi AI Service để phân tích văn bản
	log.Println("Bắt đầu gọi AI để phân tích...")
	aiResultString, err := services.AnalyzeText(textContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed: " + err.Error()})
		return
	}

	// Parse chuỗi JSON từ AI vào struct
	cleanedJSONString := cleanAIResponse(aiResultString)

	var analysisResp AnalysisResponse
	err = json.Unmarshal([]byte(cleanedJSONString), &analysisResp)
	if err != nil {
		log.Printf("Lỗi khi parse JSON từ AI: %v. \nChuỗi nhận được: %s", err, aiResultString)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response."})
		return
	}

	// Lưu kết quả vào database cùng với file hash
	analysisModel := models.Analysis{
		Summary:        analysisResp.Summary,
		KeyClauses:     analysisResp.KeyClauses,
		PotentialRisks: analysisResp.PotentialRisks,
		FileHash:       fileHash,
	}
	if err := database.DB.Create(&analysisModel).Error; err != nil {
		log.Printf("Failed to save analysis to DB: %v", err)
	}

	// Trả về kết quả có cấu trúc hoàn chỉnh cho frontend
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