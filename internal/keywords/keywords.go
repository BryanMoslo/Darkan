package keywords

import (
	validation "darkan/internal/validation"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// Keyword is a model that represents a Keyword item
// in the database
type Keyword struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Value       string    `json:"value" db:"value"`
	CallbackURL string    `json:"callback_url" db:"callback_url"`
	Found       bool      `json:"found" db:"found"`
}

// KeywordService is the interface that wraps the basic CRUD operations
// for the Keyword model
type Service interface {
	Create(keyword *Keyword) error
	UnfoundKeywords() ([]Keyword, error)
}

func (k Keyword) ValidateValue() validation.Rules {
	value := strings.TrimSpace(k.Value)
	return validation.Rules{
		func() error {
			if value == "" {
				return fmt.Errorf("%s can't be blank", "Keyword value")
			}
			return nil
		},
		func() error {
			if len(value) <= 1 {
				return fmt.Errorf("%s should contain more than 1 character", "Keyword value")
			}
			return nil
		},
	}
}

func (k Keyword) ValidateCallback() validation.Rules {
	value := strings.TrimSpace(k.CallbackURL)
	return validation.Rules{
		func() error {
			if value == "" {
				return fmt.Errorf("%s can't be blank", "Callback value")
			}
			return nil
		},
		// Here's needed the validation to check if the callback is reachable.
	}
}
