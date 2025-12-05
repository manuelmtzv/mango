package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name            string `form:"name" validate:"required,min=3,max=50"`
		Username        string `form:"username" validate:"required,min=3,max=30"`
		Email           string `form:"email" validate:"required,email"`
		Password        string `form:"password" validate:"required,strongpassword"`
		ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
	}

	if err := s.decodeForm(r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	locale := r.Context().Value(localeKey).(string)

	if err := s.validateStruct(input); err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		validationErrors := s.formatValidationErrors(err)
		var message string
		if errors, ok := validationErrors["errors"].(map[string]map[string]string); ok {
			for _, errData := range errors {
				fieldName := s.i18n.Translate(locale, fmt.Sprintf("field.%s", errData["field"]))
				message = s.i18n.Translate(locale, errData["key"], map[string]any{
					"Field": fieldName,
					"Param": errData["param"],
				})
				break
			}
		} else {
			message = s.i18n.Translate(locale, "validation.failed")
		}

		s.renderBlock(w, "alert-error", map[string]any{
			"Message": message,
		})
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
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		s.renderBlock(w, "alert-error", map[string]any{
			"Message": s.i18n.Translate(locale, "register.error.email_taken"),
		})
		return
	}

	sessionDuration := time.Duration(s.cfg.SessionDurationHours) * time.Hour
	sessionID, err := s.session.CreateSession(r.Context(), user.ID, sessionDuration)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   s.cfg.SessionDurationHours * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/dashboard", locale))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `form:"email" validate:"required"`
		Password string `form:"password" validate:"required"`
	}

	if err := s.decodeForm(r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	locale := r.Context().Value(localeKey).(string)

	if err := s.validateStruct(input); err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		validationErrors := s.formatValidationErrors(err)
		var message string
		if errors, ok := validationErrors["errors"].(map[string]map[string]string); ok {
			for _, errData := range errors {
				fieldName := s.i18n.Translate(locale, fmt.Sprintf("field.%s", errData["field"]))
				message = s.i18n.Translate(locale, errData["key"], map[string]any{
					"Field": fieldName,
					"Param": errData["param"],
				})
				break
			}
		} else {
			message = s.i18n.Translate(locale, "validation.failed")
		}

		s.renderBlock(w, "alert-error", map[string]any{
			"Message": message,
		})
		return
	}

	user, err := s.store.Users.GetByEmailOrUsername(r.Context(), input.Email)
	if err != nil || user == nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		s.renderBlock(w, "alert-error", map[string]any{
			"Message": s.i18n.Translate(locale, "login.error.invalid_credentials"),
		})
		return
	}

	match, err := VerifyPassword(input.Password, user.Hash)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if !match {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		s.renderBlock(w, "alert-error", map[string]any{
			"Message": s.i18n.Translate(locale, "login.error.invalid_credentials"),
		})
		return
	}

	sessionDuration := time.Duration(s.cfg.SessionDurationHours) * time.Hour
	sessionID, err := s.session.CreateSession(r.Context(), user.ID, sessionDuration)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   s.cfg.SessionDurationHours * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/dashboard", locale))
	w.WriteHeader(http.StatusOK)
}
