package keywords

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Match is a model that represents a match between a Keyword and a source
// in the database
type Match struct {
	ID        uuid.UUID `json:"id" db:"id"`
	KeywordID uuid.UUID `json:"keyword_id" db:"keyword_id"`
	Source    string    `json:"source_url" db:"source_url"`
	Content   string    `json:"content" db:"content"`
	FoundAt   time.Time `json:"found_at" db:"found_at"`
}

type Matches []Match

// MatchService is the interface that wraps the CRUD operations
// for the Match model
type MatchService interface {
	Create(match *Match) error
}
