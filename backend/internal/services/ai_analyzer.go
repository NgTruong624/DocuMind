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
func AnalyzeText(textContent string) (string, error) {
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

	// Initialize the model
	model := client.GenerativeModel("gemini-1.5-flash")

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