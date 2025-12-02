package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/manuelmtzv/mangocatnotes-api/internal/store"
)

type application struct {
	config Config
	store  *store.Storage
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Starting server on port %s\n", app.config.Port)
	return srv.ListenAndServe()
}
