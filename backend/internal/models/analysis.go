// backend/internal/models/analysis.go
package models

import (
	"time"

	"github.com/lib/pq"
)

type Analysis struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Summary        string         `json:"summary"`
	KeyClauses     pq.StringArray `gorm:"type:text[]" json:"key_clauses"`
	PotentialRisks pq.StringArray `gorm:"type:text[]" json:"potential_risks"`
	CreatedAt      time.Time      `json:"created_at"`
	FileHash       string         `gorm:"uniqueIndex;not null" json:"file_hash"`
}