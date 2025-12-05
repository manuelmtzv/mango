package server

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/manuelmtzv/mangocatnotes-api/internal/models"
)

func (s *Server) render(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	locale := r.Context().Value(localeKey).(string)

	t := func(key string) string {
		return s.i18n.Translate(locale, key)
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"t": t,
		"add": func(a, b int64) int64 {
			return a + b
		},
		"sub": func(a, b int64) int64 {
			return a - b
		},
	}

	partials, err := filepath.Glob("web/templates/partials/*.html")
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	fragments, err := filepath.Glob("web/templates/fragments/*.html")
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	files := append([]string{
		filepath.Join("web", "templates", "layouts", "layout.html"),
		filepath.Join("web", "templates", "pages", page),
	}, partials...)
	files = append(files, fragments...)

	tmpl, err := template.New("layout.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		s.errorJSON(w, fmt.Errorf("error parsing template: %w", err), http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = make(map[string]any)
	}

	data["Lang"] = locale
	data["t"] = t
	data["Locales"] = s.i18n.GetLocales()

	if title, ok := data["Title"].(string); ok {
		data["Title"] = s.i18n.Translate(locale, title)
	}
	if desc, ok := data["Description"].(string); ok {
		data["Description"] = s.i18n.Translate(locale, desc)
	} else {
		data["Description"] = s.i18n.Translate(locale, "meta.description")
	}

	data["CurrentPath"] = r.URL.Path
	data["BaseURL"] = s.cfg.BaseURL
	data["CurrentYear"] = time.Now().Year()
	data["AssetVersion"] = s.AssetVersion

	pathWithoutLocale := strings.TrimPrefix(r.URL.Path, "/"+locale)
	if pathWithoutLocale == "" {
		pathWithoutLocale = "/"
	}
	data["PathWithoutLocale"] = pathWithoutLocale

	err = tmpl.Execute(w, data)
	if err != nil {
		s.errorJSON(w, fmt.Errorf("error executing template: %w", err), http.StatusInternalServerError)
	}
}

func (s *Server) renderBlock(w http.ResponseWriter, r *http.Request, block string, data map[string]any) {
	locale := r.Context().Value(localeKey).(string)
	t := func(key string) string {
		return s.i18n.Translate(locale, key)
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"t": t,
	}

	partials, err := filepath.Glob("web/templates/partials/*.html")
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	fragments, err := filepath.Glob("web/templates/fragments/*.html")
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	files := append(partials, fragments...)

	tmpl, err := template.New("partials").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		s.errorJSON(w, fmt.Errorf("error parsing template: %w", err), http.StatusInternalServerError)
		return
	}

	if data == nil {
		data = make(map[string]any)
	}
	data["Lang"] = locale

	if err := tmpl.ExecuteTemplate(w, block, data); err != nil {
		s.errorJSON(w, fmt.Errorf("error executing template block: %w", err), http.StatusInternalServerError)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	features := []map[string]string{
		{
			"icon":     "notebook-pen",
			"titleKey": "features.notes.title",
			"descKey":  "features.notes.description",
		},
		{
			"icon":     "globe",
			"titleKey": "features.multilang.title",
			"descKey":  "features.multilang.description",
		},
		{
			"icon":     "tags",
			"titleKey": "features.tags.title",
			"descKey":  "features.tags.description",
		},
	}

	s.render(w, r, "index.html", map[string]any{
		"Title":       "index.title",
		"Description": "index.description",
		"Features":    features,
	})
}

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "login.html", map[string]any{
		"Title": "login.title",
	})
}

func (s *Server) registerPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "register.html", map[string]any{
		"Title": "register.title",
	})
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
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

	s.render(w, r, "dashboard.html", map[string]any{
		"Title": "dashboard.title",
		"Notes": notes,
		"Meta": models.PaginationMetadata{
			Page:       (page),
			Limit:      (limit),
			Count:      count,
			TotalPages: (count + limit - 1) / limit,
		},
	})
}
