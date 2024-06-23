package response

import (
	"encoding/json"
	"net/http"
)

// Response is generic response structure
type Response struct {
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	StatusCode int         `json:"status_code"`
}

type Responder struct {
}

func sendJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func (rs Responder) SendSuccess(w http.ResponseWriter, statusCode int, payload interface{}) {
	response := Response{
		Data:       payload,
		StatusCode: statusCode,
	}
	sendJson(w, statusCode, response)
}

func (rs Responder) SendError(w http.ResponseWriter, statusCode int, message string) {
	response := Response{
		Error:      message,
		StatusCode: statusCode,
	}
	sendJson(w, statusCode, response)
}
