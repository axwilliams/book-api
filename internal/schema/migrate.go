package schema

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var users string
	_ = tx.QueryRow("SELECT to_regclass('users')").Scan(&users)

	if users == "" {
		q := `CREATE TABLE IF NOT EXISTS users(
						id UUID,
						username varchar(255) NOT NULL,
						email varchar(255) UNIQUE NOT NULL,
						roles varchar(255)[],
						password varchar(255) NOT NULL,
						PRIMARY KEY (id)
					);`

		_, err := tx.Exec(q)
		if err != nil {
			return fmt.Errorf("Creating table: users: %w", err)
		}
	}

	var book string
	_ = tx.QueryRow("SELECT to_regclass('book')").Scan(&book)

	if book == "" {
		q := `CREATE TABLE IF NOT EXISTS book(
						id UUID, isbn varchar(255) NOT NULL,
						title varchar(255) NOT NULL,
						author varchar(255) NOT NULL,
						category varchar(255) NULL,
						PRIMARY KEY (id)
					);`

		_, err := tx.Exec(q)
		if err != nil {
			return fmt.Errorf("Creating table: book: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Committing: %w", err)
	}

	fmt.Println("Migration complete")

	return nil
}
