package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/axwilliams/books-api/cmd/books-api/handlers"
	"github.com/axwilliams/books-api/internal/business/book"
	"github.com/axwilliams/books-api/internal/business/user"
	"github.com/axwilliams/books-api/internal/middleware"
	"github.com/axwilliams/books-api/internal/platform/auth"
	"github.com/axwilliams/books-api/internal/platform/database"
	"github.com/axwilliams/books-api/internal/platform/database/postgres"
	"github.com/axwilliams/books-api/internal/schema"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log := log.New(os.Stdout, "API : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Println("[error]", err)
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

	command := flag.String("command", "", "")
	flag.Parse()

	if len(*command) != 0 {
		switch *command {
		case "migrate":
			err = schema.Migrate(db, log)
		case "seed":
			err = schema.Seed(db, log)
		default:
			return fmt.Errorf("Unknown command: %+v", err)
		}

		if err != nil {
			return fmt.Errorf("Executing %s command: %+v", *command, err)
		}
	}

	bookRepository := book.NewRepository(db)
	bookService := book.NewService(bookRepository)
	bookHandler := handlers.NewBookHandler(bookService)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	mux := mux.NewRouter()
	api := mux.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.Authenticate)
	api.Use(middleware.Logger)

	api.HandleFunc("/books", bookHandler.FindAll).Methods("GET")
	api.HandleFunc("/books/{id}", bookHandler.FindById).Methods("GET")
	api.HandleFunc("/books", middleware.HasRole(bookHandler.Add, auth.RoleAuthor)).Methods("POST")
	api.HandleFunc("/books/{id}", middleware.HasRole(bookHandler.Edit, auth.RoleAuthor)).Methods("PATCH")
	api.HandleFunc("/books/{id}", middleware.HasRole(bookHandler.Delete, auth.RoleAuthor)).Methods("DELETE")

	api.HandleFunc("/users", middleware.HasRole(userHandler.Add, auth.RoleAdmin)).Methods("POST")
	api.HandleFunc("/users/{id}", middleware.HasRole(userHandler.Edit, auth.RoleAdmin)).Methods("PATCH")
	api.HandleFunc("/users/{id}", middleware.HasRole(userHandler.Delete, auth.RoleAdmin)).Methods("DELETE")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	svrErrs := make(chan error, 1)

	go func() {
		log.Printf("[main] API listening on %s", os.Getenv("SVR_PORT"))
		svrErrs <- http.ListenAndServe(os.Getenv("SVR_PORT"), api)
	}()

	select {
	case err := <-svrErrs:
		return fmt.Errorf("Server error: %+v", err)
	case sig := <-shutdown:
		log.Printf("[main] Starting shutdown: %v", sig)

		// TODO: Add shutdown code
	}

	return nil
}
