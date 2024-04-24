package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

type Middleware func(next http.Handler) http.Handler

func Create(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			x := ms[i]
			next = x(next)
		}
		return next
	}
}

// Logging middleware logs relevant information about the request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		start := time.Now()
		next.ServeHTTP(ww, r)
		log.Println(ww.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

// IsAuthenticated middleware checks authentication against a secret stored in secret manager
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretToken := os.Getenv("API_KEY_SECRET")
		if secretToken == "" {
			log.Panic("Secret not set")
		}

		// Retrieve the Authorization header.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Trim prefix and check match
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token.
		if bearerToken != secretToken {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, continue serving the request.
		next.ServeHTTP(w, r)
	})
}
