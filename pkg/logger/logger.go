package logger

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func New() *Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	logger, _ := cfg.Build()
	return &Logger{logger.Sugar()}
}

func NewNop() *Logger {
	logger := zap.NewNop()
	return &Logger{logger.Sugar()}
}

func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		l.Infow("request", "method", r.Method, "path", r.URL.Path, "requestId", reqID)
		next.ServeHTTP(w, r)
	})
}
