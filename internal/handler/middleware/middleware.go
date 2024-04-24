package middleware

import (
	"log"
	"net/http"
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
