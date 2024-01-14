package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetUser(w, r)
	case "POST":
		return s.handleCreateUser(w, r)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	default:
		return fmt.Errorf("not allowed %s methods", r.Method)
	}

}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	if id, ok := mux.Vars(r)["id"]; ok {
		// db := db.Get user id
		fmt.Println(id)
		return ApiJSONResponse(w, http.StatusFound, &User{})
	}

	user := NewTestUser("admin", "example@admin.pl", "Haslo#123")
	log.Println(r)
	return ApiJSONResponse(w, http.StatusFound, user)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	return ApiJSONResponse(w, http.StatusFound, "")
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	return ApiJSONResponse(w, http.StatusFound, "")
}

type ApiResponse struct {
	TotalCount int         `json:"total_count"`
	Result     interface{} `json:"result"`
	PageSize   int         `json:"page_size"`
	Page       int         `json:"page"`
}

func ApiJSONResponse(w http.ResponseWriter, status int, value interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := &ApiResponse{
		Result: value,
	}
	return json.NewEncoder(w).Encode(*response)
}

type apiFunc func(http.ResponseWriter, *http.Request) error
type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			ApiJSONResponse(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func runAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/user/{id}", makeHTTPHandleFunc(s.handleGetUser))

	log.Println("JSON Api server running on port: ", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
	}
}