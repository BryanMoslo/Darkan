package internal

import (
	"github.com/leapkit/core/db"
	"github.com/leapkit/core/envor"
	_ "github.com/lib/pq"
)

var (
	// DatabaseURL is the connection string for the database
	// that will be used by the application.
	DatabaseURL = envor.Get("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/darkan_development?sslmode=disable")

	// DB is the database connection builder function
	// that will be used by the application based on the driver and
	// connection string.
	DB = db.ConnectionFn(DatabaseURL, db.WithDriver("postgres"))
)
