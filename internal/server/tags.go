package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) createTag(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	var input models.Tag
	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	existingTags, _, err := s.store.Tags.GetAll(r.Context(), userID, 1, 1000)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for _, tag := range existingTags {
		if tag.Name == input.Name {
			s.errorJSON(w, fmt.Errorf("Tag with provided name (%s) already exists.", input.Name), http.StatusBadRequest)
			return
		}
	}

	input.UserID = userID
	if err := s.store.Tags.Create(r.Context(), &input); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusCreated, input)
}

func (s *Server) getTags(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	tags, count, err := s.store.Tags.GetAll(r.Context(), userID, 1, 1000)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	s.writeJSON(w, http.StatusOK, map[string]any{
		"data":  tags,
		"count": count,
	})
}

func (s *Server) getTag(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	tag, err := s.store.Tags.GetByID(r.Context(), id)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if tag == nil {
		s.errorJSON(w, errors.New("tag not found"), http.StatusNotFound)
		return
	}

	s.writeJSON(w, http.StatusOK, tag)
}

func (s *Server) updateTag(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var input models.Tag
	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	input.ID = id
	if err := s.store.Tags.Update(r.Context(), &input); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, input)
}

func (s *Server) deleteTag(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	if err := s.store.Tags.Delete(r.Context(), id); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) findTagsOrCreate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	var input struct {
		Tags []string `json:"tags"`
	}
	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	tags, err := s.store.Tags.FindOrCreate(r.Context(), userID, input.Tags)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]any{
		"data":  tags,
		"count": len(tags),
	})
}
