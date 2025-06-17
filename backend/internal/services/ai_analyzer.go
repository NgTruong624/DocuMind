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
		Analyze the following contract text and return the analysis as a single JSON string.
		IMPORTANT: Your response must ONLY contain the JSON string, with no additional text, explanations, or markdown formatting.

		The JSON must follow this exact structure:
		{
			"summary": "A concise, professional summary of the contract's main points",
			"key_clauses": ["List of the most important clauses as an array of strings"],
			"potential_risks": ["List of potential risks or points of concern as an array of strings. Return empty array [] if none found"]
		}

		Contract text to analyze:
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