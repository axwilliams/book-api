package schema

import (
	"database/sql"
	"fmt"
	"log"
)

func Migrate(db *sql.DB, log *log.Logger) error {
	log.Println("[schema] Starting migration")

	q := `CREATE TABLE IF NOT EXISTS users(
					id UUID,
					username varchar(255) NOT NULL,
					email varchar(255) UNIQUE NOT NULL,
					roles varchar(255)[],
					password varchar(255) NOT NULL,
					PRIMARY KEY (id)
				);`

	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("Creating table: users: %w", err)
	}

	q = `CREATE TABLE IF NOT EXISTS book(
					id UUID, isbn varchar(255) NOT NULL,
					title varchar(255) NOT NULL,
					author varchar(255) NOT NULL,
					category varchar(255) NULL,
					PRIMARY KEY (id)
				);`

	_, err = db.Exec(q)
	if err != nil {
		return fmt.Errorf("Creating table: book: %w", err)
	}

	log.Println("[schema] Migration complete")

	return nil
}
