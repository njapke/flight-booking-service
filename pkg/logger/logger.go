package logger

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

type Logger struct {
	*log.Logger
}

func New() *Logger {
	return &Logger{log.New(os.Stderr, "", log.LstdFlags)}
}

func (l *Logger) Error(err error) {
	l.Printf("ERROR: %s", err)
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := middleware.GetReqID(r.Context())
		l.Printf("[%s] %s %s", reqId, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
