package handler

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Result bool        `json:"result"`
	Data   interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Result bool     `json:"result"`
	Errors []string `json:"errors"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteSuccess writes a success response
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, SuccessResponse{
		Result: true,
		Data:   data,
	})
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, statusCode int, errors ...string) {
	WriteJSON(w, statusCode, ErrorResponse{
		Result: false,
		Errors: errors,
	})
}
