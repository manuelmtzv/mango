package server

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) getMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	user, err := s.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, user)
}
