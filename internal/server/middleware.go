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

func (s *Server) GuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.Context().Value(localeKey)
		if locale == nil {
			locale = "es"
		}

		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		_, err = s.session.GetSession(r.Context(), cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/dashboard", locale), http.StatusFound)
	})
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
