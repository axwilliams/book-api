package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/axwilliams/books-api/internal/platform/auth"
	"github.com/axwilliams/books-api/internal/platform/web"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidID      = errors.New("ID is not in the correct form")
	ErrUsernameExists = errors.New("username is already taken")
	ErrEmailExists    = errors.New("email is already taken")
)

type Service interface {
	GetById(id string) (*User, error)
	Create(nu *NewUser) (*User, error)
	Update(id string, uu UpdateUser) error
	Destroy(id string) error
	Authenticate(username, password string) (auth.Claims, error)
}

type service struct {
	ur Repository
}

func NewService(ur Repository) Service {
	return &service{
		ur,
	}
}

func (s *service) GetById(id string) (*User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	return s.ur.GetById(id)
}

func (s *service) Create(nu *NewUser) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("generating password hash: %w", err)
	}

	u := &User{
		ID:           uuid.New().String(),
		Username:     strings.TrimSpace(nu.Username),
		Email:        strings.TrimSpace(nu.Email),
		Roles:        nu.Roles,
		PasswordHash: hash,
	}

	if ok := s.ur.UsernameAvailable(u.Username, ""); !ok {
		return nil, web.NewRequestError(ErrUsernameExists, http.StatusNotAcceptable)
	}

	if ok := s.ur.EmailAvailable(u.Email, ""); !ok {
		return nil, web.NewRequestError(ErrEmailExists, http.StatusNotAcceptable)
	}

	return u, s.ur.Create(u)
}

func (s *service) Update(id string, uu UpdateUser) error {
	if _, err := uuid.Parse(id); err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	u, err := s.ur.GetById(id)
	switch {
	case err == ErrNoUserFound:
		return web.NewRequestError(ErrNoAffect, http.StatusGone)
	case err != nil:
		return err
	}

	if uu.Username != "" {
		u.Username = strings.TrimSpace(uu.Username)
		if ok := s.ur.UsernameAvailable(u.Username, u.ID); !ok {
			return web.NewRequestError(ErrUsernameExists, http.StatusNotAcceptable)
		}
	}

	if uu.Email != "" {
		u.Email = strings.TrimSpace(uu.Email)
		if ok := s.ur.EmailAvailable(u.Email, u.ID); !ok {
			return web.NewRequestError(ErrEmailExists, http.StatusNotAcceptable)
		}
	}

	if len(uu.Roles) != 0 {
		u.Roles = uu.Roles
	}

	if uu.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		u.PasswordHash = hash
	}

	return s.ur.Update(u)
}

func (s *service) Destroy(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	return s.ur.Destroy(id)
}

func (s *service) Authenticate(username, password string) (auth.Claims, error) {
	u, err := s.ur.GetByUsername(username)
	switch {
	case err == ErrNoUserFound:
		return auth.Claims{}, web.NewRequestError(ErrBadCredentials, http.StatusUnauthorized)
	case err != nil:
		return auth.Claims{}, err
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, web.NewRequestError(ErrBadCredentials, http.StatusUnauthorized)
	}

	claims := auth.NewClaims(u.ID, u.Roles)

	return claims, nil
}
