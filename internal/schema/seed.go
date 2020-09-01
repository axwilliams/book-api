package schema

import (
	"database/sql"
	"fmt"
	"log"
)

func Seed(db *sql.DB, log *log.Logger) error {
	log.Println("[schema] Seeding data")

	q := `INSERT INTO users (id, username, email, password, roles)VALUES (
					'a72bec75-0a5f-49af-a844-5763d188788e', 
					'admin', 
					'admin@example.com', 
					'$2a$10$IUs9j88n5g5pthZXNmU9tei2mhIX7MWTvk39AjWUx40juWOrrPOzi', 
					'{ADMIN}'
				),(
					'69a47775-6d89-4d38-ad38-acdb2928f6a1', 
					'author', 
					'author@example.com', 
					'$2a$10$ExnMCA7MuOwW.s8Ss0BvSuGNCHawMIpqMmyJ4Oa9sTCTKKw2x445e', 
					'{AUTHOR}'
				),(
					'bad069ce-4afa-4a53-a673-14ae7b627d06', 
					'user', 
					'user@example.com', 
					'$2a$10$PH3juwOeFvwr0auAgrUOy.PeCOGT/oRZevj7I5urM7tcOhsw0bdIi', 
					'{}' 
				);`

	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("Seeding table: users: %w", err)
	}

	q = `INSERT INTO book (id, isbn, title, author, category) VALUES
					('f4ac7e14-fc8e-4096-b956-34e5a33040f2', '978-0241372579', 'The Castle', 'Franz Kafka', 'Fiction'),
					('71432eb9-58da-4eae-aa20-ccc49064246f', '978-1451673319', 'Fahrenheit 451', 'Ray Bradbury', 'Fiction'),
					('562e1fe0-0dde-4717-a008-cd2a699301d2', '978-0465025275', 'Six Easy Pieces', 'Richard Feynman', 'Science'
				);`

	_, err = db.Exec(q)
	if err != nil {
		return fmt.Errorf("[error] seeding table: book: %w", err)
	}

	log.Println("[schema] Seeding complete")

	return nil
}
