package middleware

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"path"
)

var sLvl int

func init() {
	sLvl = 0
}

func clean(p string) string {
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
