package server

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const localeKey contextKey = "locale"

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
	filesDir := http.Dir(filepath.Join(workDir, "web/static"))
	FileServer(r, "/static", filesDir)

	r.Get("/", s.handleRoot)

	r.Route("/{locale}", func(r chi.Router) {
		r.Use(s.localeMiddleware)

		r.Get("/", s.home)
		r.Get("/login", s.loginPage)
		r.Get("/register", s.registerPage)

		r.Group(func(r chi.Router) {
			r.Use(s.AuthMiddleware)
			r.Get("/dashboard", s.dashboard)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", s.register)
			r.Post("/login", s.login)
			r.With(s.AuthMiddleware).Post("/logout", s.logout)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(s.AuthMiddleware)
			r.Get("/me", s.getMe)
		})

		r.Route("/notes", func(r chi.Router) {
			r.Use(s.AuthMiddleware)
			r.Get("/", s.getNotes)
			r.Get("/new", s.createNotePage)
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
	})

	return r
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/es/", http.StatusFound)
}

func (s *Server) localeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := chi.URLParam(r, "locale")

		if locale != "en" && locale != "es" && locale != "it" {
			http.NotFound(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), localeKey, locale)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(
			path,
			http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP,
		)
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
