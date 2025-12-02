package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,min=3,max=30"`
		Password string `json:"password" validate:"required,strongpassword"`
		Name     string `json:"name"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := app.validateStruct(input); err != nil {
		app.writeJSON(w, http.StatusBadRequest, app.formatValidationErrors(err))
		return
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:    input.Email,
		Username: input.Username,
		Hash:     hash,
		Name:     input.Name,
	}

	if err := app.store.Users.Create(r.Context(), user); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken(user.ID.Hex(), app.config.JWTSecret)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, models.LoginResponse{
		Username:    user.Username,
		AccessToken: token,
	})
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Identifier string `json:"identifier" validate:"required"`
		Password   string `json:"password" validate:"required"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := app.validateStruct(input); err != nil {
		app.writeJSON(w, http.StatusBadRequest, app.formatValidationErrors(err))
		return
	}

	user, err := app.store.Users.GetByEmailOrUsername(r.Context(), input.Identifier)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	if user == nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	match, err := VerifyPassword(input.Password, user.Hash)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if !match {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user.ID.Hex(), app.config.JWTSecret)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, models.LoginResponse{
		Username:    user.Username,
		AccessToken: token,
	})
}

func (app *application) getMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(primitive.ObjectID)
	user, err := app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, user)
}

func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(primitive.ObjectID)
	var input models.Note
	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	input.UserID = userID
	if err := app.store.Notes.Create(r.Context(), &input); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), input.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	input.Tags = tags

	app.writeJSON(w, http.StatusCreated, input)
}

func (app *application) getNotes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(primitive.ObjectID)
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

	notes, count, err := app.store.Notes.GetAll(r.Context(), userID, int64(page), int64(limit), search, tags)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}

	for i := range notes {
		noteTags, err := app.store.Notes.GetTags(r.Context(), notes[i].ID)
		if err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		notes[i].Tags = noteTags
	}

	app.writeJSON(w, http.StatusOK, models.PaginatedNotesResponse{
		Data: notes,
		Meta: models.PaginationMetadata{
			Page:       (page),
			Limit:      (limit),
			Count:      count,
			TotalPages: (count + limit - 1) / limit,
		},
	})
}

func (app *application) getNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	note, err := app.store.Notes.GetByID(r.Context(), id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if note == nil {
		app.errorJSON(w, errors.New("note not found"), http.StatusNotFound)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), note.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	app.writeJSON(w, http.StatusOK, note)
}

func (app *application) updateNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var input models.Note
	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	input.ID = id
	if err := app.store.Notes.Update(r.Context(), &input); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), input.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	input.Tags = tags

	app.writeJSON(w, http.StatusOK, input)
}

func (app *application) deleteNote(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	if err := app.store.Notes.Delete(r.Context(), id); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) createTag(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(primitive.ObjectID)
	var input models.Tag
	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	existingTags, _, err := app.store.Tags.GetAll(r.Context(), userID, 1, 1000)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for _, tag := range existingTags {
		if tag.Name == input.Name {
			app.errorJSON(w, errors.New(fmt.Sprintf("Tag with provided name (%s) already exists.", input.Name)), http.StatusBadRequest)
			return
		}
	}

	input.UserID = userID
	if err := app.store.Tags.Create(r.Context(), &input); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, input)
}

func (app *application) getTags(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(primitive.ObjectID)

	tags, count, err := app.store.Tags.GetAll(r.Context(), userID, 1, 1000) // Legacy fetches all, setting high limit
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  tags,
		"count": count,
	})
}

func (app *application) getTag(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	tag, err := app.store.Tags.GetByID(r.Context(), id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if tag == nil {
		app.errorJSON(w, errors.New("tag not found"), http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, tag)
}

func (app *application) updateTag(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	var input models.Tag
	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	input.ID = id
	if err := app.store.Tags.Update(r.Context(), &input); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, input)
}

func (app *application) deleteTag(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"), http.StatusBadRequest)
		return
	}

	if err := app.store.Tags.Delete(r.Context(), id); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) findTagsOrCreate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(primitive.ObjectID)
	var input struct {
		Tags []string `json:"tags"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	tags, err := app.store.Tags.FindOrCreate(r.Context(), userID, input.Tags)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  tags,
		"count": len(tags),
	})
}

func (app *application) getNoteTags(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := primitive.ObjectIDFromHex(noteIDParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, tags)
}

func (app *application) attachNoteTags(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := primitive.ObjectIDFromHex(noteIDParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	var input struct {
		Tags []string `json:"tags"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	var newTagIDs []primitive.ObjectID
	for _, idStr := range input.Tags {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			app.errorJSON(w, errors.New("invalid tag id in list"), http.StatusBadRequest)
			return
		}
		newTagIDs = append(newTagIDs, id)
	}

	existingTags, err := app.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
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
			if err := app.store.Notes.DetachTag(r.Context(), noteID, existingTag.ID); err != nil {
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
		}
	}

	if len(newTagIDs) > 0 {
		if err := app.store.Notes.AttachTags(r.Context(), noteID, newTagIDs); err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
	}

	note, err := app.store.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	app.writeJSON(w, http.StatusOK, note)
}

func (app *application) attachNoteTag(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := primitive.ObjectIDFromHex(noteIDParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	tagIDParam := chi.URLParam(r, "tagId")
	tagID, err := primitive.ObjectIDFromHex(tagIDParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid tag id"), http.StatusBadRequest)
		return
	}

	if err := app.store.Notes.AttachTags(r.Context(), noteID, []primitive.ObjectID{tagID}); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	note, err := app.store.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	app.writeJSON(w, http.StatusOK, note)
}

func (app *application) detachNoteTag(w http.ResponseWriter, r *http.Request) {
	noteIDParam := chi.URLParam(r, "id")
	noteID, err := primitive.ObjectIDFromHex(noteIDParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid note id"), http.StatusBadRequest)
		return
	}

	tagIDParam := chi.URLParam(r, "tagId")
	tagID, err := primitive.ObjectIDFromHex(tagIDParam)
	if err != nil {
		app.errorJSON(w, errors.New("invalid tag id"), http.StatusBadRequest)
		return
	}

	if err := app.store.Notes.DetachTag(r.Context(), noteID, tagID); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	note, err := app.store.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	tags, err := app.store.Notes.GetTags(r.Context(), noteID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	note.Tags = tags

	app.writeJSON(w, http.StatusOK, note)
}
