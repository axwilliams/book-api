package test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/axwilliams/books-api/internal/platform/database"
	"github.com/axwilliams/books-api/internal/platform/database/postgres"
	"github.com/axwilliams/books-api/internal/schema"
	pg_test "github.com/axwilliams/books-api/internal/test/postgres"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func Setup() (*sql.DB, *pg_test.Container) {
	log := log.New(os.Stdout, "TEST: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	container := pg_test.LaunchContainer(log)

	db, err := postgres.Open(database.Config{
		User:     "postgres",
		Password: "postgres",
		Host:     container.Host,
		Name:     "postgres",
	})
	if err != nil {
		log.Fatalf("[error] Connecting to test database: %v", err)
	}

	fmt.Println("Waiting for database ...")

	var pingError error
	for attempts := 1; attempts <= 20; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		Teardown(db, container)
		log.Fatalf("Database never ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		Teardown(db, container)
		log.Fatalf("[error] Creating test tables: %v", err)
	}

	if err := schema.Seed(db); err != nil {
		Teardown(db, container)
		log.Fatalf("[error] Seeding test data: %v", err)
	}

	return db, container
}

func Teardown(db *sql.DB, container *pg_test.Container) {
	db.Close()
	pg_test.DestroyContainer(container)
}
