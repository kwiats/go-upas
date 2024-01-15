package models

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIError struct {
	Error string
}
type ApiResponse struct {
	TotalCount int         `json:"total_count,omitempty"`
	Result     interface{} `json:"result"`
	PageSize   int         `json:"page_size,omitempty"`
	Page       int         `json:"page,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, value interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// if value is list then calc totalCount , pageSize and page
	if status == http.StatusNotFound {
		log.Print(value)
	}
	response := &ApiResponse{
		Result: value,
	}
	return json.NewEncoder(w).Encode(*response)
}
