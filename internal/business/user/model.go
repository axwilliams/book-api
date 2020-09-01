package user

import (
	"github.com/lib/pq"
)

type User struct {
	ID           string         `db:"id" json:"id"`
	Username     string         `db:"username" json:"username"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"roles"`
	PasswordHash []byte         `db:"password" json:"-"`
}

type NewUser struct {
	Username string   `json:"username" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Roles    []string `json:"roles"`
	Password string   `json:"password" validate:"required,password"`
}

type UpdateUser struct {
	Username string   `json:"username"`
	Email    string   `json:"email" validate:"omitempty,email"`
	Roles    []string `json:"roles"`
	Password string   `json:"password" validate:"omitempty,password"`
}
