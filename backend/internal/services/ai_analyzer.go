package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AnalyzeText sends text to Google Gemini for analysis and returns a structured JSON response
func AnalyzeText(textContent string, modelName ...string) (string, error) {
	// Step 1: Initialize the client
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error creating Gemini client: %v", err)
		return "", fmt.Errorf("failed to initialize Gemini client: %w", err)
	}
	defer client.Close()

	// Smart model selection based on content length
	chosenModel := selectOptimalModel(textContent)
	if len(modelName) > 0 && modelName[0] != "" {
		chosenModel = modelName[0]
	}
	
	log.Printf("Using model: %s for content length: %d characters", chosenModel, len(textContent))
	model := client.GenerativeModel(chosenModel)

	// Step 2: Construct the prompt
	prompt := fmt.Sprintf(`
	Phân tích nội dung hợp đồng sau và trả về kết quả bằng tiếng Việt dưới dạng một chuỗi JSON duy nhất.
	QUAN TRỌNG: Phản hồi của bạn CHỈ ĐƯỢC chứa chuỗi JSON, không có văn bản, giải thích hay định dạng markdown nào khác.

	JSON phải tuân theo cấu trúc chính xác sau:
	{
		"summary": "Một bản tóm tắt chuyên nghiệp, ngắn gọn bằng tiếng Việt về các điểm chính của hợp đồng",
		"key_clauses": ["Danh sách các điều khoản quan trọng nhất bằng tiếng Việt, dưới dạng một mảng các chuỗi"],
		"potential_risks": ["Danh sách các rủi ro tiềm ẩn hoặc các điểm cần lưu ý bằng tiếng Việt, dưới dạng một mảng các chuỗi. Trả về mảng rỗng [] nếu không tìm thấy"]
	}

	Nội dung hợp đồng cần phân tích:
	---
	%s
	---
`, textContent)

	// Step 3: Send request to AI
	log.Println("Sending request to Gemini API...")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		
		// Kiểm tra lỗi quota API
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "exceeded") {
			return "", fmt.Errorf("API quota exceeded: %w", err)
		}
		
		// Kiểm tra lỗi API key
		if strings.Contains(err.Error(), "API_KEY") || strings.Contains(err.Error(), "authentication") {
			return "", fmt.Errorf("API authentication failed: %w", err)
		}
		
		// Lỗi chung
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Step 4: Process and return the response
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received empty response from Gemini API")
	}

	if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		// Clean the response string to ensure it's valid JSON
		cleanedResponse := strings.TrimSpace(string(txt))
		if cleanedResponse == "" {
			return "", fmt.Errorf("received empty analysis from Gemini API")
		}
		return cleanedResponse, nil
	}

	return "", fmt.Errorf("failed to extract text from Gemini API response")
}

// selectOptimalModel chọn model tối ưu dựa trên độ dài nội dung
func selectOptimalModel(content string) string {
	contentLength := len(content)
	
	// Kiểm tra các tiêu chí để chọn model
	if shouldUsePro(content, contentLength) {
		log.Printf("Selected Gemini Pro 2.5 - Content analysis: length=%d chars", contentLength)
		return GeminiPro25
	}
	
	log.Printf("Selected Gemini Flash 2.5 - Content analysis: length=%d chars", contentLength)
	return GeminiFlash25
}

// shouldUsePro quyết định có nên dùng model Pro không
func shouldUsePro(content string, length int) bool {
	// Tiêu chí 1: Độ dài nội dung
	if length > ModelSwitchThreshold {
		return true
	}
	
	// Tiêu chí 2: Độ phức tạp dựa trên từ khóa pháp lý
	complexLegalKeywords := []string{
		"bồi thường", "vi phạm", "tranh chấp", "kiện tụng", "phạt", 
		"lãi suất", "thế chấp", "bảo lãnh", "trách nhiệm pháp lý",
		"điều khoản phạt", "force majeure", "bất khả kháng",
		"quyền sở hữu trí tuệ", "bản quyền", "thương hiệu",
		"miễn trừ trách nhiệm", "hủy bỏ hợp đồng", "chấm dứt",
	}
	
	contentLower := strings.ToLower(content)
	complexKeywordCount := 0
	
	for _, keyword := range complexLegalKeywords {
		if strings.Contains(contentLower, keyword) {
			complexKeywordCount++
		}
	}
	
	// Nếu có >= 3 từ khóa phức tạp, dùng Pro
	if complexKeywordCount >= 3 {
		log.Printf("Complex legal content detected (%d keywords), switching to Pro", complexKeywordCount)
		return true
	}
	
	// Tiêu chí 3: Số lượng điều khoản (đếm "điều", "khoản", "mục")
	clauseKeywords := []string{"điều ", "khoản ", "mục ", "chương "}
	clauseCount := 0
	
	for _, keyword := range clauseKeywords {
		clauseCount += strings.Count(contentLower, keyword)
	}
	
	// Nếu có >= 10 điều khoản, dùng Pro
	if clauseCount >= 10 {
		log.Printf("Multiple clauses detected (%d), switching to Pro", clauseCount)
		return true
	}
	
	return false
}

// AnalyzeTextSmart - Wrapper function với tự động chọn model thông minh
func AnalyzeTextSmart(textContent string) (string, error) {
	return AnalyzeText(textContent) // Sẽ tự động chọn model qua selectOptimalModel
}

// AskContractQuestionSmart - Wrapper function với tự động chọn model thông minh  
func AskContractQuestionSmart(contractText, question string) (string, error) {
	return AskContractQuestion(contractText, question) // Sẽ tự động chọn model qua selectOptimalModel
}

func AskContractQuestion(contractText, question string, modelName ...string) (string, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to initialize Gemini client: %w", err)
	}
	defer client.Close()

	// Smart model selection based on content length
	chosenModel := selectOptimalModel(contractText)
	if len(modelName) > 0 && modelName[0] != "" {
		chosenModel = modelName[0]
	}
	
	log.Printf("Using model: %s for contract length: %d characters", chosenModel, len(contractText))
	model := client.GenerativeModel(chosenModel)

	prompt := fmt.Sprintf(`
Bạn là một trợ lý pháp lý. Dựa trên nội dung hợp đồng sau, hãy trả lời NGẮN GỌN, rõ ràng, bằng tiếng Việt cho câu hỏi của người dùng. 
Chỉ trả lời nội dung liên quan, không cần giải thích thêm, không trả về JSON, chỉ trả lời như hội thoại tự nhiên.

Nội dung hợp đồng:
---
%s
---

Câu hỏi: %s
`, contractText, question)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		// Kiểm tra lỗi quota API
		if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "exceeded") {
			return "", fmt.Errorf("API quota exceeded: %w", err)
		}
		
		// Kiểm tra lỗi API key
		if strings.Contains(err.Error(), "API_KEY") || strings.Contains(err.Error(), "authentication") {
			return "", fmt.Errorf("API authentication failed: %w", err)
		}
		
		// Lỗi chung
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received empty response from Gemini API")
	}

	if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return strings.TrimSpace(string(txt)), nil
	}

	return "", fmt.Errorf("failed to extract text from Gemini API response")
}

// Helper function để tạo constants cho các model names
const (
	GeminiFlash25 = "gemini-2.5-flash"
	GeminiPro25   = "gemini-2.5-pro"
	// Ngưỡng ký tự để chuyển từ Flash sang Pro
	ModelSwitchThreshold = 15000 // 15k ký tự
)

// Convenience functions để sử dụng các model cụ thể
func AnalyzeTextWithFlash25(textContent string) (string, error) {
	return AnalyzeText(textContent, GeminiFlash25)
}

func AnalyzeTextWithPro25(textContent string) (string, error) {
	return AnalyzeText(textContent, GeminiPro25)
}

func AskContractQuestionWithFlash25(contractText, question string) (string, error) {
	return AskContractQuestion(contractText, question, GeminiFlash25)
}

func AskContractQuestionWithPro25(contractText, question string) (string, error) {
	return AskContractQuestion(contractText, question, GeminiPro25)
}