package database

import (
	"fmt"
	"os"

	"documind/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// THAY ĐỔI: Thêm models.AnalysisDetail{} vào AutoMigrate
	// GORM sẽ tự động tạo cả hai bảng `analyses` và `analysis_details`
	if err := db.AutoMigrate(&models.Analysis{}, &models.AnalysisDetail{}); err != nil {
		return nil, fmt.Errorf("auto-migration failed: %w", err)
	}

	DB = db
	return db, nil
}