package handlers

import (
	"net/http"

	"github.com/axwilliams/book-api/internal/business/user"
	"github.com/axwilliams/book-api/internal/platform/auth"
	"github.com/axwilliams/book-api/internal/platform/web"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	us user.Service
}

func NewUserHandler(us user.Service) UserHandler {
	return UserHandler{
		us,
	}
}

func (h *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	nu := user.NewUser{}

	if err := web.Decode(r, &nu); err != nil {
		web.RespondError(w, err)
		return
	}

	u, err := h.us.Create(&nu)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, web.Message("id", u.ID), http.StatusCreated)
}

func (h *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uu := user.UpdateUser{}
	if err := web.Decode(r, &uu); err != nil {
		web.RespondError(w, err)
		return
	}

	if err := h.us.Update(vars["id"], uu); err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, nil, http.StatusOK)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.us.Destroy(vars["id"]); err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, nil, http.StatusOK)
}

func (h *UserHandler) Token(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		web.RespondError(w, web.NewRequestError(user.ErrBasicAuth, http.StatusUnauthorized))
		return
	}

	claims, err := h.us.Authenticate(username, password)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	token, err := auth.CreateToken(claims)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.Respond(w, web.Message("token", token), http.StatusOK)
}
