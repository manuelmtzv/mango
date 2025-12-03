package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) routes() http.Handler {
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

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	FileServer(r, "/assets", filesDir)

	r.Get("/", s.home)
	r.Get("/login", s.loginPage)
	r.Get("/register", s.registerPage)

	r.Group(func(r chi.Router) {

		r.Get("/dashboard", s.dashboard)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.register)
		r.Post("/login", s.login)
		r.With(s.AuthMiddleware).Get("/validate-token", func(w http.ResponseWriter, r *http.Request) {
			token := GetJwt(r)
			userID := r.Context().Value(UserIDKey).(uuid.UUID)
			user, err := s.store.Users.GetByID(r.Context(), userID)
			if err != nil {
				s.errorJSON(w, err, http.StatusInternalServerError)
				return
			}
			s.writeJSON(w, http.StatusOK, models.LoginResponse{
				Username:    user.Username,
				AccessToken: token,
			})
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(s.AuthMiddleware)
		r.Get("/me", s.getMe)
	})

	r.Route("/notes", func(r chi.Router) {
		r.Use(s.AuthMiddleware)
		r.Get("/", s.getNotes)
		r.Post("/", s.createNote)
		r.Get("/{id}", s.getNote)
		r.Patch("/{id}", s.updateNote)
		r.Delete("/{id}", s.deleteNote)

		r.Get("/{id}/tags", s.getNoteTags)
		r.Patch("/{id}/tags", s.attachNoteTags)
		r.Patch("/{id}/tags/{tagId}", s.attachNoteTag)
		r.Delete("/{id}/tags/{tagId}", s.detachNoteTag)
	})

	r.Route("/tags", func(r chi.Router) {
		r.Use(s.AuthMiddleware)
		r.Get("/", s.getTags)
		r.Post("/", s.createTag)
		r.Post("/find-or-create", s.findTagsOrCreate)
		r.Get("/{id}", s.getTag)
		r.Patch("/{id}", s.updateTag)
		r.Delete("/{id}", s.deleteTag)
	})

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
