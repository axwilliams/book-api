package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axwilliams/books-api/platform/database"
	"github.com/axwilliams/books-api/platform/database/postgres"
	"github.com/joho/godotenv"
)

func main() {
	log := log.New(os.Stdout, "API : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("[error]: ", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	if err := godotenv.Load("../../.env"); err != nil {
		return fmt.Errorf("Loading env file: %+v", err)
	}

	db, err := postgres.Open(database.Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
	})
	if err != nil {
		return fmt.Errorf("Connecting to database: %+v", err)
	}

	defer func() {
		log.Println("[main] Stopping database")
		db.Close()
	}()

	return nil
}
