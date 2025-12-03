package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	var input models.Note
	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	input.UserID = userID
	if err := s.store.Notes.Create(r.Context(), &input); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), input.ID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	input.Tags = tags

	s.writeJSON(w, http.StatusCreated, input)
}

func (s *Server) getNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit < 1 {
		limit = 10
	}

	search := r.URL.Query().Get("search")
	tags := r.URL.Query()["tags"]

	notes, count, err := s.store.Notes.GetAll(r.Context(), userID, int64(page), int64(limit), search, tags)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}

	for i := range notes {
		noteTags, err := s.store.Notes.GetTags(r.Context(), notes[i].ID)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		notes[i].Tags = noteTags
	}

	s.writeJSON(w, http.StatusOK, models.PaginatedNotesResponse{
		Data: notes,
		Meta: models.PaginationMetadata{
			Page:       (page),
			Limit:      (limit),
			Count:      count,
			TotalPages: (count + limit - 1) / limit,
		},
	})
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	note, err := s.store.Notes.GetByID(r.Context(), id)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if note == nil {
		s.errorJSON(w, errors.New("note not found"), http.StatusNotFound)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), note.ID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	s.writeJSON(w, http.StatusOK, note)
}

func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var input models.Note
	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	input.ID = id
	if err := s.store.Notes.Update(r.Context(), &input); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), input.ID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	input.Tags = tags

	s.writeJSON(w, http.StatusOK, input)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	if err := s.store.Notes.Delete(r.Context(), id); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) getNoteTags(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := uuid.Parse(noteIDParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, http.StatusOK, tags)
}

func (s *Server) attachNoteTags(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := uuid.Parse(noteIDParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	var input struct {
		Tags []string `json:"tags"`
	}
	if err := s.readJSON(w, r, &input); err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	var newTagIDs []uuid.UUID
	for _, idStr := range input.Tags {
		id, err := uuid.Parse(idStr)
		if err != nil {
			s.errorJSON(w, errors.New("invalid tag id in list"), http.StatusBadRequest)
			return
		}
		newTagIDs = append(newTagIDs, id)
	}

	existingTags, err := s.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for _, existingTag := range existingTags {
		found := false
		for _, newTagID := range newTagIDs {
			if existingTag.ID == newTagID {
				found = true
				break
			}
		}
		if !found {
			if err := s.store.Notes.DetachTag(r.Context(), noteID, existingTag.ID); err != nil {
				s.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
		}
	}

	if len(newTagIDs) > 0 {
		if err := s.store.Notes.AttachTags(r.Context(), noteID, newTagIDs); err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	note, err := s.store.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	s.writeJSON(w, http.StatusOK, note)
}

func (s *Server) attachNoteTag(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := uuid.Parse(noteIDParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	tagIDParam := chi.URLParam(r, "tagId")
	tagID, err := uuid.Parse(tagIDParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid tag id"), http.StatusBadRequest)
		return
	}

	if err := s.store.Notes.AttachTags(r.Context(), noteID, []uuid.UUID{tagID}); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	note, err := s.store.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	s.writeJSON(w, http.StatusOK, note)
}

func (s *Server) detachNoteTag(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := uuid.Parse(noteIDParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	tagIDParam := chi.URLParam(r, "tagId")
	tagID, err := uuid.Parse(tagIDParam)
	if err != nil {
		s.errorJSON(w, errors.New("invalid tag id"), http.StatusBadRequest)
		return
	}

	if err := s.store.Notes.DetachTag(r.Context(), noteID, tagID); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	note, err := s.store.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	s.writeJSON(w, http.StatusOK, note)
}
