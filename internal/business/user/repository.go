package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/axwilliams/books-api/internal/platform/web"
)

var (
	ErrNoAffect       = errors.New("No rows affected")
	ErrBadCredentials = errors.New("Invalid username or password")
	ErrNoUserFound    = errors.New("No user found")
)

type Repository interface {
	GetById(id string) (*User, error)
	Create(u *User) error
	Update(u *User) error
	Destroy(id string) error
	GetByUsername(username string) (*User, error)
	UsernameAvailable(username, cuurentID string) bool
	EmailAvailable(username, cuurentID string) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db,
	}
}

func (r *repository) GetById(id string) (*User, error) {
	u := &User{}

	err := r.db.QueryRow("SELECT id, username, email, roles, password FROM users where id=$1",
		id).Scan(&u.ID, &u.Username, &u.Email, &u.Roles, &u.PasswordHash)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNoUserFound
	case err != nil:
		return nil, fmt.Errorf("Retrieving user: %w", err)
	}

	return u, nil
}

func (r *repository) Create(u *User) error {
	_, err := r.db.Exec("INSERT INTO users (id, username, email, roles, password) VALUES ($1, $2, $3, $4, $5)",
		u.ID, u.Username, u.Email, u.Roles, u.PasswordHash)

	if err != nil {
		return fmt.Errorf("Creating user: %w", err)
	}

	return nil
}

func (r *repository) Update(u *User) error {
	_, err := r.db.Exec("UPDATE users SET username=$1, email=$2, roles=$3, password=$4 WHERE id=$5;",
		u.Username, u.Email, u.Roles, u.PasswordHash, u.ID)

	if err != nil {
		return fmt.Errorf("Updating user: %w", err)
	}

	return nil
}

func (r *repository) Destroy(id string) error {
	res, err := r.db.Exec("DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Counting affected users: %w", err)
	}

	if count <= 0 {
		return web.NewRequestError(ErrNoAffect, http.StatusGone)
	}

	return nil
}

func (r *repository) GetByUsername(username string) (*User, error) {
	u := &User{}

	err := r.db.QueryRow("SELECT id, username, email, roles, password FROM users WHERE username = $1",
		username).Scan(&u.ID, &u.Username, &u.Email, &u.Roles, &u.PasswordHash)

	switch {
	case err == sql.ErrNoRows:
		return nil, ErrNoUserFound
	case err != nil:
		return nil, fmt.Errorf("Retrieving user by username: %w", err)
	}

	return u, nil
}

func (r *repository) UsernameAvailable(username, currentID string) bool {

	var ID string

	err := r.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&ID)

	if err == sql.ErrNoRows || ID == currentID {
		return true
	}

	return false

}

func (r *repository) EmailAvailable(email, currentID string) bool {
	var ID string

	err := r.db.QueryRow("SELECT id FROM users WHERE email = $1 ", email).Scan(&ID)

	if err == sql.ErrNoRows || ID == currentID {
		return true
	}

	return false
}
