package main

import (
	"log"
	"net/http"
	"os"

	"documind/backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Tải các biến môi trường từ file configs/.env
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Lấy PORT từ biến môi trường
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port mặc định nếu không được thiết lập
	}
	
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Đăng ký route cho việc phân tích file
	r.POST("/api/v1/analyze", handlers.AnalyzeHandler) 

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}