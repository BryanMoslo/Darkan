package keywords

import (
	"darkan/internal/sources"
	"strings"

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

func (s *service) UnfoundKeywords() ([]Instance, error) {
	var keywords []Instance
	err := s.db.Select(&keywords, "SELECT * FROM keywords WHERE found = $1", false)
	return keywords, err
}

func (s *service) CreateMatch(match *Match) error {
	match.ID = uuid.Must(uuid.NewV4())
	_, err := s.db.NamedExec(`INSERT INTO matches (id, keyword_id, source_url, content, found_at) VALUES (:id, :keyword_id, :source_url, :content, :found_at)`, match)
	return err
}

func (s *service) SourceList() (sources sources.Sources, err error) {
	err = s.db.Select(&sources, "SELECT * FROM sources")

	return sources, err
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint")
}
