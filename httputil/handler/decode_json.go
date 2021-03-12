package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" && r.Header.Get("Content-Type") != "" {
		return NewAPIError(http.StatusUnsupportedMediaType, "content_type.unsupported")
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB.

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return NewAPIError(http.StatusBadRequest, "body.invalid_json")

		case errors.Is(err, io.ErrUnexpectedEOF):
			return NewAPIError(http.StatusBadRequest, "body.invalid_json")

		case errors.As(err, &unmarshalTypeError):
			return NewAPIError(http.StatusBadRequest, fmt.Sprintf("%s.invalid_type", unmarshalTypeError.Field))

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return NewAPIError(http.StatusBadRequest, fmt.Sprintf("%s.unknown_field", strings.Trim(fieldName, "\"")))

		case errors.Is(err, io.EOF):
			return NewAPIError(http.StatusBadRequest, "body.empty")

		case err.Error() == "http: request body too large":
			return NewAPIError(http.StatusRequestEntityTooLarge, "body.too_large")

		default:
			return err
		}
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return NewAPIError(http.StatusBadRequest, "body.invalid_json")
	}

	return nil
}
