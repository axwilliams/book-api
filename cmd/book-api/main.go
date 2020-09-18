package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axwilliams/book-api/cmd/book-api/handlers"
	"github.com/axwilliams/book-api/internal/business/book"
	"github.com/axwilliams/book-api/internal/business/user"
	"github.com/axwilliams/book-api/internal/middleware"
	"github.com/axwilliams/book-api/internal/platform/auth"
	"github.com/axwilliams/book-api/internal/platform/database"
	"github.com/axwilliams/book-api/internal/platform/database/postgres"
	"github.com/axwilliams/book-api/internal/platform/web"
	"github.com/axwilliams/book-api/internal/schema"
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

	fmt.Println("Waiting for database ...")

	var pingError error
	for attempts := 1; attempts <= 20; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if err := schema.Migrate(db); err != nil {
		return fmt.Errorf("Migrating tables: %+v", err)
	}

	if err := schema.Seed(db); err != nil {
		return fmt.Errorf("Seeding data: %+v", err)
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

	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.Respond(w, "Books API v1", http.StatusOK)
	})

	api.HandleFunc("/books", bookHandler.FindAll).Methods("GET")
	api.HandleFunc("/books/{id}", bookHandler.FindById).Methods("GET")
	api.HandleFunc("/search/books", bookHandler.Search).Methods("GET")
	api.HandleFunc("/books", middleware.HasRole(bookHandler.Add, auth.RoleAuthor)).Methods("POST")
	api.HandleFunc("/books/{id}", middleware.HasRole(bookHandler.Edit, auth.RoleAuthor)).Methods("PATCH")
	api.HandleFunc("/books/{id}", middleware.HasRole(bookHandler.Delete, auth.RoleAuthor)).Methods("DELETE")

	api.HandleFunc("/users", middleware.HasRole(userHandler.Add, auth.RoleAdmin)).Methods("POST")
	api.HandleFunc("/users/{id}", middleware.HasRole(userHandler.Edit, auth.RoleAdmin)).Methods("PATCH")
	api.HandleFunc("/users/{id}", middleware.HasRole(userHandler.Delete, auth.RoleAdmin)).Methods("DELETE")

	api.HandleFunc("/users/token", userHandler.Token).Methods("POST")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	svrErrs := make(chan error, 1)

	srv := http.Server{
		Addr:         os.Getenv("SVR_PORT"),
		Handler:      api,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		log.Printf("[main] API listening on %s", srv.Addr)
		svrErrs <- srv.ListenAndServe()
	}()

	select {
	case err := <-svrErrs:
		return fmt.Errorf("Server error: %+v", err)

	case sig := <-shutdown:
		log.Printf("[main] Starting graceful shutdown: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			srv.Close()
			return fmt.Errorf("Graceful shutdown failed: %+v", err)
		}
	}

	return nil
}
