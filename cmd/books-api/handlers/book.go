package handlers

import (
	"net/http"
	"strings"

	"github.com/axwilliams/books-api/internal/business/book"
	"github.com/axwilliams/books-api/internal/platform/web"
	"github.com/gorilla/mux"
)

type BookHandler struct {
	bs book.Service
}

func NewBookHandler(bs book.Service) BookHandler {
	return BookHandler{
		bs,
	}
}

func (h *BookHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	bks, err := h.bs.GetAll()
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, bks, http.StatusOK)
}

func (h *BookHandler) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bk, err := h.bs.GetById(vars["id"])
	if err != nil && err != book.ErrNoBookFound {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, bk, http.StatusOK)
}

func (h *BookHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	params := book.SearchParams{}
	params.ISBN = strings.TrimSpace(q.Get("isbn"))
	params.Title = strings.TrimSpace(q.Get("title"))
	params.Author = strings.TrimSpace(q.Get("author"))
	params.Category = strings.TrimSpace(q.Get("category"))

	sort := strings.TrimSpace(q.Get("sort"))
	order := strings.TrimSpace(q.Get("order"))
	limitStr := strings.TrimSpace(q.Get("limit"))
	offsetStr := strings.TrimSpace(q.Get("offset"))

	bks, err := h.bs.Search(params, sort, order, limitStr, offsetStr)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, bks, http.StatusOK)
}

func (h *BookHandler) Add(w http.ResponseWriter, r *http.Request) {
	nb := book.NewBook{}

	if err := web.Decode(r, &nb); err != nil {
		web.RespondError(w, err)
		return
	}

	bk, err := h.bs.Create(&nb)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, web.Message("id", bk.ID), http.StatusCreated)
}

func (h *BookHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ub := book.UpdateBook{}
	if err := web.Decode(r, &ub); err != nil {
		web.RespondError(w, err)
		return
	}

	if err := h.bs.Update(vars["id"], ub); err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, nil, http.StatusOK)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.bs.Destroy(vars["id"]); err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, nil, http.StatusOK)
}
