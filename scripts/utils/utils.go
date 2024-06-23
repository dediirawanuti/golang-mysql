package utils

import (
	"encoding/json"
	"net/http"
	// "github.com/golang-mysql/scripts/connection"
)

type Response struct {
	Status int64 `json:"status"`
	Message string `json:"message"`
	Error string `json:"error"`
}

func Healthz(w http.ResponseWriter, r *http.Request) {
		var response Response
		response.Status = 0
		response.Message = "OK"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
}
