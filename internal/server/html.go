package server

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (s *Server) render(w http.ResponseWriter, _ *http.Request, page string, data any) {
	files := []string{
		filepath.Join("web", "templates", "layouts", "base.html"),
		filepath.Join("web", "templates", "pages", page),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "index.html", map[string]any{
		"Title": "Home",
	})
}

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "login.html", map[string]any{
		"Title": "Login",
	})
}

func (s *Server) registerPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "register.html", map[string]any{
		"Title": "Register",
	})
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "dashboard.html", map[string]any{
		"Title": "Dashboard",
	})
}
