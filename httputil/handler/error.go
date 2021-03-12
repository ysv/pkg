package handler

import (
	"encoding/json"
)

// TODO: Add multiple errors handling.
type APIError struct {
	Status int    `json:"-"`
	Name   string `json:"name"`
}

// Satisfy the error interface.
func (apiErr APIError) Error() string {
	return apiErr.Name
}

// Follow internal errors format convention.
func (apiErr APIError) MarshalJSON() ([]byte, error) {
	type Alert struct {
		Name string `json:"name"`
	}
	res := struct {
		Alerts []Alert `json:"alerts"`
	}{
		Alerts: []Alert{{Name: apiErr.Name}},
	}

	return json.Marshal(&res)
}

func NewAPIError(status int, name string) *APIError {
	return &APIError{
		Status: status,
		Name:   name,
	}
}
