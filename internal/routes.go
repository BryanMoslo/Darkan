package internal

import (
	"darkan/internal/keywords"

	"github.com/leapkit/core/server"
)

// AddRoutes mounts the routes for the application,
// it assumes that the base services have been injected
// in the creation of the server instance.
func AddRoutes(r server.Router) error {
	// TODO: Authentication

	r.HandleFunc("POST /api/search", keywords.Create)

	return nil
}
