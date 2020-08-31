package book

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/axwilliams/books-api/internal/platform/web"
)

var (
	ErrNoAffect    = errors.New("No rows affected")
	ErrNoBookFound = errors.New("No book found")
)

type Repository interface {
	GetAll() ([]Book, error)
	GetById(id string) (*Book, error)
	Create(bk *Book) error
	Update(bk *Book) error
	Destroy(id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db,
	}
}

func (r *repository) GetAll() ([]Book, error) {
	rows, err := r.db.Query("SELECT id, isbn, title, author, category FROM book")
	if err != nil {
		return nil, fmt.Errorf("Retrieving books: %w", err)
	}
	defer rows.Close()

	bks := []Book{}
	for rows.Next() {
		bk := Book{}
		if err = rows.Scan(&bk.ID, &bk.ISBN, &bk.Title, &bk.Author, &bk.Category); err != nil {
			return nil, fmt.Errorf("Scanning book rows: %w", err)
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterating book rows: %w", err)
	}

	return bks, nil
}

func (r *repository) GetById(id string) (*Book, error) {
	bk := &Book{}

	err := r.db.QueryRow("SELECT id, isbn, title, author, category FROM book where id=$1",
		id).Scan(&bk.ID, &bk.ISBN, &bk.Title, &bk.Author, &bk.Category)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNoBookFound
	case err != nil:
		return nil, fmt.Errorf("Retrieving book: %w", err)
	}

	return bk, nil
}

func (r *repository) Create(bk *Book) error {
	_, err := r.db.Exec("INSERT INTO book (id, isbn, title, author, category) VALUES ($1, $2, $3, $4, $5)",
		bk.ID, bk.ISBN, bk.Title, bk.Author, bk.Category)

	if err != nil {
		return fmt.Errorf("Creating book: %w", err)
	}

	return nil
}

func (r *repository) Update(bk *Book) error {
	_, err := r.db.Exec("UPDATE book SET isbn=$1, title=$2, author=$3, category=$4 WHERE id=$5;",
		bk.ISBN, bk.Title, bk.Author, bk.Category, bk.ID)

	if err != nil {
		return fmt.Errorf("Updating book: %w", err)
	}

	return nil
}

func (r *repository) Destroy(id string) error {
	res, err := r.db.Exec("DELETE FROM book WHERE id = $1;", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Counting affected books: %w", err)
	}

	if count <= 0 {
		return web.NewRequestError(ErrNoAffect, http.StatusGone)
	}

	return nil
}
