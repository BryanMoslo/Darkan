package sources

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

// Source is a model that represents a dark web source in the database
type Source struct {
	ID        uuid.UUID `json:"id" db:"id"`
	URL       string    `json:"url" db:"url"`
	Available bool      `json:"available" db:"available"`
	RiskLevel int       `json:"risk_level" db:"risk_level"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Sources []Source
