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
	// Chỉ tải file .env nếu không ở trong môi trường production ("release")
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load("configs/.env")
		if err != nil {
			// Thay vì làm sập chương trình, chỉ in ra một cảnh báo.
			log.Println("Warning: .env file not found, using system environment variables.")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	r := gin.Default()

	// Khởi tạo kết nối database trong một goroutine để không chặn việc khởi động server.
	// Đây là chìa khóa để khắc phục lỗi "Timed Out" trên Render.
	go func() {
		log.Println("Attempting to connect to the database...")
		_, dbErr := database.Connect()
		if dbErr != nil {
			log.Printf("!!! Asynchronous database connection failed: %v", dbErr)
		} else {
			log.Println("Database connection established successfully.")
		}
	}()

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

	// Endpoint health check cho Render
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Các API endpoints của ứng dụng
	api := r.Group("/api/v1")
	{
		api.POST("/analyze", handlers.AnalyzeHandler)
		api.POST("/contract-chat", handlers.ContractChatHandler)
		api.GET("/analyses", handlers.GetAnalyses)
		api.GET("/analyses/:id", handlers.GetAnalysisDetail)
	}

	log.Printf("Server starting on port %s", port)
	// Dòng này sẽ khởi động server ngay lập tức
	r.Run(":" + port)
}