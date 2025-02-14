package router

import (
	"encoding/json"
	"net/http"
)

type Errors struct {
	Error string `json:"errors"`
}

func WriteErrorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Errors{err.Error()})
}
