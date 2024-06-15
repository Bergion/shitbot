package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func jsonResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	responseBody, _ := json.Marshal(response)
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}

func Ok(w http.ResponseWriter, data interface{}) {
	jsonResponse(w, http.StatusOK,
		Response{
			Success: true,
			Data:    data,
		},
	)
}

func NotFound(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusNotFound,
		Response{
			Success: false,
			Message: message,
		},
	)
}

func BadRequest(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusBadRequest,
		Response{
			Success: false,
			Message: message,
		},
	)
}

func Forbidden(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusForbidden,
		Response{
			Success: false,
			Message: message,
		},
	)
}

func InternalServerError(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusInternalServerError,
		Response{
			Success: false,
			Message: message,
		},
	)
}
