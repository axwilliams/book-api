package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/axwilliams/books-api/internal/platform/auth"
	"github.com/axwilliams/books-api/internal/platform/web"
)

var (
	ErrAuthHeader = errors.New(("Wrong authorization header format"))
	ErrDenied     = errors.New(("Permission denied"))
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/api/v1/users/token" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			web.RespondError(w, web.NewRequestError(ErrAuthHeader, http.StatusBadRequest))
			return
		}

		claims, err := auth.ParseWithClaims(parts[1])
		if err != nil {
			web.RespondError(w, err)
			return
		}

		ctx := auth.ContextWithUser(r.Context(), claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HasRole(next http.HandlerFunc, role string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		roles, ok := auth.RolesFromContext(r.Context())
		if !ok {
			web.RespondError(w, web.NewRequestError(ErrDenied, http.StatusForbidden))
			return
		}

		if !auth.HasRole(roles, role) {
			web.RespondError(w, web.NewRequestError(ErrDenied, http.StatusForbidden))
			return
		}

		next(w, r)
	}
}
