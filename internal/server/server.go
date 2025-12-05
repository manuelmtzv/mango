package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/manuelmtzv/mangocatnotes-api/internal/config"
	"github.com/manuelmtzv/mangocatnotes-api/internal/i18n"
	"github.com/manuelmtzv/mangocatnotes-api/internal/session"
	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
)

type Server struct {
	cfg          *config.Config
	i18n         *i18n.Manager
	store        *store.Storage
	session      *session.SessionManager
	AssetVersion string
}

func New(cfg *config.Config, store *store.Storage, session *session.SessionManager) *Server {
	return &Server{
		cfg:          cfg,
		i18n:         i18n.NewManager(),
		store:        store,
		session:      session,
		AssetVersion: fmt.Sprintf("%d", time.Now().Unix()),
	}
}

func (s *Server) Start() error {
	if err := s.i18n.LoadDir("web/locales"); err != nil {
		return fmt.Errorf("failed to load translations: %w", err)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.cfg.Port),
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Starting server on port %s\n", s.cfg.Port)
	return srv.ListenAndServe()
}
