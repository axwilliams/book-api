package handlers

import (
	"net/http"

	"github.com/axwilliams/books-api/internal/business/book"
	"github.com/axwilliams/books-api/internal/platform/web"
	"github.com/gorilla/mux"
)

type BookHandler interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	FindById(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type bookHandler struct {
	bs book.Service
}

func NewBookHandler(bs book.Service) BookHandler {
	return &bookHandler{
		bs,
	}
}

func (h *bookHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	bks, err := h.bs.GetAll()
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, bks, http.StatusOK)
}

func (h *bookHandler) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bk, err := h.bs.GetById(vars["id"])
	if err != nil && err != book.ErrNoBookFound {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, bk, http.StatusOK)
}

func (h *bookHandler) Add(w http.ResponseWriter, r *http.Request) {
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

func (h *bookHandler) Edit(w http.ResponseWriter, r *http.Request) {
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

func (h *bookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.bs.Destroy(vars["id"]); err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, nil, http.StatusOK)
}
