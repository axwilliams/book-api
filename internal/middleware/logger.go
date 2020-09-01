package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/axwilliams/books-api/internal/platform/auth"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *loggingResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		UserID, ok := auth.UserFromContext(r.Context())
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		rw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		log.Printf("[request] %s : (%d) : %s %s -> %s (%s)",
			UserID,
			rw.statusCode,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
