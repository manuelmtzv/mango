package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=3,max=30"`
		Password string `json:"password" validate:"required,strongpassword"`
		Name     string `json:"name"`
	}

	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := s.validateStruct(input); err != nil {
		s.writeJSON(w, http.StatusBadRequest, s.formatValidationErrors(err))
		return
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:    input.Email,
		Username: input.Username,
		Hash:     hash,
		Name:     input.Name,
	}

	if err := s.store.Users.Create(r.Context(), user); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	sessionID, err := s.cache.CreateSession(r.Context(), user.ID, 7*24*time.Hour)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	s.writeJSON(w, http.StatusCreated, map[string]any{
		"username": user.Username,
	})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Identifier string `json:"identifier" validate:"required"`
		Password   string `json:"password" validate:"required"`
	}

	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := s.validateStruct(input); err != nil {
		s.writeJSON(w, http.StatusBadRequest, s.formatValidationErrors(err))
		return
	}

	user, err := s.store.Users.GetByEmailOrUsername(r.Context(), input.Identifier)
	if err != nil {
		s.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	if user == nil {
		s.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	match, err := VerifyPassword(input.Password, user.Hash)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if !match {
		s.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	sessionID, err := s.cache.CreateSession(r.Context(), user.ID, 7*24*time.Hour)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	s.writeJSON(w, http.StatusOK, map[string]any{
		"username": user.Username,
	})
}
