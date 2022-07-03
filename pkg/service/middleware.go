package service

import "net/http"

func (s *Service) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					s.sendError(w, err.Error(), http.StatusInternalServerError)
				} else {
					s.log.Printf("panic: %v", r)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
