package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorResponse struct {
	Message string `json:"error"`
}

func (e errorResponse) Error() string {
	return e.Message
}

type statusResponse struct {
	Status string `json:"status"`
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		slog.Error("Failed to write header:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	slog.Error(message)

	response := errorResponse{
		Message: message,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	w.WriteHeader(statusCode)
	if err != nil {
		slog.Error("Failed to write header:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
