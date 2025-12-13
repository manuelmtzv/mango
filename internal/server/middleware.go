package server

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userID"

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.Context().Value(localeKey)
		if locale == nil {
			locale = "es"
		}

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/%s/login", locale), http.StatusFound)
			return
		}

		userID, err := s.session.GetSession(r.Context(), cookie.Value)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/%s/login", locale), http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
