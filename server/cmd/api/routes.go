package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
	"github.com/google/uuid"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.register)
		r.Post("/login", app.login)
		r.With(app.AuthMiddleware).Get("/validate-token", func(w http.ResponseWriter, r *http.Request) {
			token := GetJwt(r)
			userID := r.Context().Value(UserIDKey).(uuid.UUID)
			user, err := app.store.Users.GetByID(r.Context(), userID)
			if err != nil {
				app.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			app.writeJSON(w, http.StatusOK, models.LoginResponse{
				Username:    user.Username,
				AccessToken: token,
			})
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(app.AuthMiddleware)
		r.Get("/me", app.getMe)
	})

	r.Route("/notes", func(r chi.Router) {
		r.Use(app.AuthMiddleware)
		r.Get("/", app.getNotes)
		r.Post("/", app.createNote)
		r.Get("/{id}", app.getNote)
		r.Patch("/{id}", app.updateNote)
		r.Delete("/{id}", app.deleteNote)

		r.Get("/{id}/tags", app.getNoteTags)
		r.Patch("/{id}/tags", app.attachNoteTags)
		r.Patch("/{id}/tags/{tagId}", app.attachNoteTag)
		r.Delete("/{id}/tags/{tagId}", app.detachNoteTag)
	})

	r.Route("/tags", func(r chi.Router) {
		r.Use(app.AuthMiddleware)
		r.Get("/", app.getTags)
		r.Post("/", app.createTag)
		r.Post("/find-or-create", app.findTagsOrCreate)
		r.Get("/{id}", app.getTag)
		r.Patch("/{id}", app.updateTag)
		r.Delete("/{id}", app.deleteTag)
	})

	return r
}
