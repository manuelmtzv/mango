package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) createNotePage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)

	userTags, _, err := s.store.Tags.GetAll(r.Context(), userID, 1, 1000)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.renderBlock(w, r, "create_note_modal", map[string]any{
		"AvailableTags": userTags,
	})
}

func (s *Server) editNotePage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
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

	if note.UserID != userID {
		s.errorJSON(w, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	noteTags, err := s.store.Notes.GetTags(r.Context(), note.ID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = noteTags

	userTags, _, err := s.store.Tags.GetAll(r.Context(), note.UserID, 1, 1000)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	s.renderBlock(w, r, "edit_note_modal", map[string]any{
		"Note":          note,
		"AvailableTags": userTags,
	})
}

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
	var input models.Note
	var tagNames []string

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := s.readJSON(w, r, &input); err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		input.Title = r.FormValue("title")
		input.Content = r.FormValue("content")

		tagsJSON := r.FormValue("tags")
		if tagsJSON != "" && tagsJSON != "[]" {
			if err := json.Unmarshal([]byte(tagsJSON), &tagNames); err != nil {
				s.errorJSON(w, err, http.StatusBadRequest)
				return
			}
		}
	}

	input.UserID = userID
	if err := s.store.Notes.Create(r.Context(), &input); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if len(tagNames) > 0 {
		createdTags, err := s.store.Tags.FindOrCreate(r.Context(), userID, tagNames)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		var tagIDs []uuid.UUID
		for _, tag := range createdTags {
			tagIDs = append(tagIDs, tag.ID)
		}

		if err := s.store.Notes.AttachTags(r.Context(), input.ID, tagIDs); err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	tags, err := s.store.Notes.GetTags(r.Context(), input.ID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	input.Tags = tags

	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Trigger", "note-created")
		data := map[string]any{
			"Title":     input.Title,
			"Content":   input.Content,
			"ID":        input.ID,
			"UpdatedAt": input.UpdatedAt,
			"Tags":      input.Tags,
		}
		s.renderBlock(w, r, "note-card", data)
		w.Write([]byte(`<div id="empty-state" hx-swap-oob="delete"></div>`))
		w.Write([]byte(`<div id="notes-grid" hx-swap-oob="removeClass:hidden"></div>`))
		return
	}

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
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
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

	if note.UserID != userID {
		s.errorJSON(w, errors.New("forbidden"), http.StatusForbidden)
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
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
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

	if note.UserID != userID {
		s.errorJSON(w, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	var input models.Note
	var tagNames []string

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := s.readJSON(w, r, &input); err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			s.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		input.Title = r.FormValue("title")
		input.Content = r.FormValue("content")

		tagsJSON := r.FormValue("tags")
		if tagsJSON != "" && tagsJSON != "[]" {
			if err := json.Unmarshal([]byte(tagsJSON), &tagNames); err != nil {
				s.errorJSON(w, err, http.StatusBadRequest)
				return
			}
		}
	}

	input.ID = id
	if err := s.store.Notes.Update(r.Context(), &input); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err := s.store.Notes.ClearTags(r.Context(), id); err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if len(tagNames) > 0 {
		createdTags, err := s.store.Tags.FindOrCreate(r.Context(), note.UserID, tagNames)
		if err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}

		var tagIDs []uuid.UUID
		for _, tag := range createdTags {
			tagIDs = append(tagIDs, tag.ID)
		}

		if err := s.store.Notes.AttachTags(r.Context(), id, tagIDs); err != nil {
			s.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	updatedNote, err := s.store.Notes.GetByID(r.Context(), id)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := s.store.Notes.GetTags(r.Context(), updatedNote.ID)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	updatedNote.Tags = tags

	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Trigger", "note-updated")
		data := map[string]any{
			"Title":     updatedNote.Title,
			"Content":   updatedNote.Content,
			"ID":        updatedNote.ID,
			"UpdatedAt": updatedNote.UpdatedAt,
			"Tags":      updatedNote.Tags,
		}
		s.renderBlock(w, r, "note-card", data)
		return
	}

	s.writeJSON(w, http.StatusOK, updatedNote)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(uuid.UUID)
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

	if note.UserID != userID {
		s.errorJSON(w, errors.New("forbidden"), http.StatusForbidden)
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
