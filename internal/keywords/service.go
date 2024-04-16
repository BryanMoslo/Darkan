package keywords

import (
	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
)

var _ Service = (*service)(nil)

type service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *service {
	return &service{db: db}
}

func (s *service) Create(keyword *Instance) error {
	keyword.ID = uuid.Must(uuid.NewV4())
	_, err := s.db.NamedExec(`INSERT INTO keywords (id, value, callback_url) VALUES (:id, :value, :callback_url)`, keyword)
	return err
}

func (s *service) All() ([]Instance, error) {
	var keywords []Instance
	err := s.db.Select(&keywords, "SELECT * FROM keywords WHERE found = $1", false)
	return keywords, err
}
