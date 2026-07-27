package apierrors

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, code, message, requestID string) {
	WriteJSON(w, status, errorEnvelope{Error: errorBody{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	}})
}
