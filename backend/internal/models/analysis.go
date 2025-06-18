// backend/internal/models/analysis.go
package models

import (
	"time"

	"github.com/lib/pq"
)

// Analysis là model cho bảng chính, chứa các thông tin nhẹ.
type Analysis struct {
	ID             uint      `gorm:"primaryKey"`
	FileHash       string    `gorm:"type:varchar(64);uniqueIndex"`
	CreatedAt      time.Time
	SummaryPreview string    `gorm:"type:varchar(200)"` // Lưu 200 ký tự đầu của summary

	// GORM relation: Một Analysis sẽ có một AnalysisDetail
	AnalysisDetail   AnalysisDetail `gorm:"foreignKey:AnalysisID"`
}

// AnalysisDetail chứa các dữ liệu văn bản dài.
type AnalysisDetail struct {
	ID             uint           `gorm:"primaryKey"`
	AnalysisID     uint           `gorm:"not null;index"` // Khóa ngoại liên kết ngược lại với bảng Analysis
	Summary        string         `gorm:"type:text"`
	KeyClauses     pq.StringArray `gorm:"type:text[]"`
	PotentialRisks pq.StringArray `gorm:"type:text[]"`
}