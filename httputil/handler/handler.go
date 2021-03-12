package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ysv/pkg/logger"
)

// The Handler helps to handle errors in one place.
type Handler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h(w, r); err != nil {
		switch e := err.(type) {
		case *APIError:
			logger.Warnf("HTTP %d - %s", e.Status, e)

			// We can retrieve the status here and write out a specific HTTP status code.
			w.WriteHeader(e.Status)
			if err := json.NewEncoder(w).Encode(e); err != nil {
				panic(err)
			}
		default:
			logger.Warnf("HTTP 500 - %s", e)

			w.WriteHeader(http.StatusInternalServerError)
			// Fixme: Use separate struct for 500.
			e = NewAPIError(http.StatusInternalServerError, "server.internal_error")
			if err := json.NewEncoder(w).Encode(e); err != nil {
				panic(err)
			}
		}
	}
}
