package main

import (
	"log"
	"net/http"
	"os"

	"documind/backend/internal/handlers"
	"documind/backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	_, dbErr := database.Connect()
	if dbErr != nil {
		log.Fatalf("Database connection failed: %v", dbErr)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/v1/analyze", handlers.AnalyzeHandler)

	r.POST("/api/v1/contract-chat", handlers.ContractChatHandler)

	r.GET("/api/v1/analyses", handlers.GetAnalyses)
	r.GET("/api/v1/analyses/:id", handlers.GetAnalysisDetail)

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}