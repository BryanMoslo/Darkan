package keywords

import "github.com/gofrs/uuid/v5"

// Instance is a model that represents a Keyword item
// in the database
type Instance struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Value       string    `json:"value" db:"value"`
	CallbackURL string    `json:"callback_url" db:"callback_url"`
	Found       bool      `json:"found" db:"found"`
}

// KeywordService is the interface that wraps the basic CRUD operations
// for the Keyword model
type Service interface {
	Create(keyword *Instance) error
}
