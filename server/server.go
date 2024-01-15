package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kwiats/go-upas/models"
	"github.com/kwiats/go-upas/storage"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if writeErr := models.WriteJSONResponse(w, http.StatusBadRequest, models.APIError{Error: err.Error()}); writeErr != nil {
				log.Printf("Error writing JSON response: %v", writeErr)
			}
		}
	}
}
func RunAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/user/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleUserWithId)))

	log.Println("JSON Api server running on port: ", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
	}
}
