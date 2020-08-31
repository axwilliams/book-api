package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/axwilliams/books-api/internal/platform/database"
	"github.com/axwilliams/books-api/internal/platform/database/postgres"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	log := log.New(os.Stdout, "API : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("[error] ", err)
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

	// bookRepository := book.NewRepository(db)
	// bookService := book.NewService(bookRepository)
	// bookHandler := handlers.NewBookHandler(bookService)

	mux := mux.NewRouter()
	api := mux.PathPrefix("/api/v1").Subrouter()

	// api.HandleFunc("/books", bookHandler.FindAll).Methods("GET")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	svrErrs := make(chan error, 1)

	go func() {
		log.Printf("[main] API listening on %s", os.Getenv("SVR_PORT"))
		svrErrs <- http.ListenAndServe(os.Getenv("SVR_PORT"), api)
	}()

	select {
	case err := <-svrErrs:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		log.Printf("[main] Starting shutdown: %v", sig)

		// TODO: Add shutdown code
	}

	return nil
}
