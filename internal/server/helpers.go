package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/schema"
)

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("form")
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type nestErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
}

func (s *Server) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) decodeForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(dst, r.Form); err != nil {
		return errors.New("error parsing form data")
	}

	return nil
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	return err
}

func (s *Server) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload nestErrorResponse
	payload.StatusCode = statusCode
	payload.Message = err.Error()
	payload.Error = http.StatusText(statusCode)

	return s.writeJSON(w, statusCode, payload)
}
