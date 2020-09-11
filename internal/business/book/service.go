package book

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/axwilliams/book-api/internal/platform/web"
	"github.com/google/uuid"
)

var (
	ErrInvalidID   = errors.New("ID is not in the correct form")
	ErrInvalidSort = errors.New("invalid sort field")
)

type Service interface {
	GetAll() ([]Book, error)
	GetById(id string) (*Book, error)
	Search(sp SearchParams, sort, order, limitStr, offsetStr string) ([]Book, error)
	Create(nb *NewBook) (*Book, error)
	Update(id string, ub UpdateBook) error
	Destroy(id string) error
}

type service struct {
	br Repository
}

func NewService(br Repository) Service {
	return &service{
		br,
	}
}

type SearchParams struct {
	ISBN     string
	Title    string
	Author   string
	Category string
}

func (s *service) GetAll() ([]Book, error) {
	return s.br.GetAll()
}

func (s *service) GetById(id string) (*Book, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	return s.br.GetById(id)
}

func (s *service) Search(sp SearchParams, sort, order, limitStr, offsetStr string) ([]Book, error) {
	sort = strings.ToLower(sort)
	order = strings.ToLower(order)

	sortOrder := ""
	if sort == "id" || sort == "isbn" || sort == "title" || sort == "author" {
		if order == "asc" || order == "desc" {
			sortOrder = sort + " " + order
		}
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 0
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	return s.br.Search(sp, sortOrder, limit, offset)
}

func (s *service) Create(nb *NewBook) (*Book, error) {
	bk := &Book{
		ID:       uuid.New().String(),
		ISBN:     strings.TrimSpace(nb.ISBN),
		Title:    strings.TrimSpace(nb.Title),
		Author:   strings.TrimSpace(nb.Author),
		Category: strings.TrimSpace(nb.Category),
	}

	return bk, s.br.Create(bk)
}

func (s *service) Update(id string, ub UpdateBook) error {
	if _, err := uuid.Parse(id); err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	bk, err := s.br.GetById(id)
	switch {
	case err == ErrNoBookFound:
		return web.NewRequestError(ErrNoAffect, http.StatusGone)
	case err != nil:
		return err
	}

	if ub.ISBN != "" {
		bk.ISBN = strings.TrimSpace(ub.ISBN)
	}
	if ub.Title != "" {
		bk.Title = strings.TrimSpace(ub.Title)
	}
	if ub.Author != "" {
		bk.Author = strings.TrimSpace(ub.Author)
	}
	if ub.Category != nil {
		bk.Category = strings.TrimSpace(*ub.Category)
	}

	return s.br.Update(bk)
}

func (s *service) Destroy(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return web.NewRequestError(ErrInvalidID, http.StatusBadRequest)
	}

	return s.br.Destroy(id)
}
