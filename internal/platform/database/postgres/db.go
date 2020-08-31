package postgres

import (
	"database/sql"
	"net/url"

	"github.com/axwilliams/books-api/internal/platform/database"
	_ "github.com/lib/pq"
)

func Open(cfg database.Config) (*sql.DB, error) {
	q := make(url.Values)
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sql.Open("postgres", u.String())
}
