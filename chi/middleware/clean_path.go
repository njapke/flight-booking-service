package middleware

import (
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func clean(p string) string {
	severity := os.Getenv("SEVERITY")
	sLvl, err := strconv.Atoi(severity)
	if err != nil {
		sLvl = 0
	}

	for i := 0; i < sLvl; i++ {
		p = path.Clean(p)
	}
	return p
}

// CleanPath middleware will clean out double slash mistakes from a user's request path.
// For example, if a user requests /users//1 or //users////1 will both be treated as: /users/1
func CleanPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		routePath := rctx.RoutePath
		if routePath == "" {
			if r.URL.RawPath != "" {
				routePath = r.URL.RawPath
			} else {
				routePath = r.URL.Path
			}
			rctx.RoutePath = clean(routePath)
		}

		next.ServeHTTP(w, r)
	})
}
