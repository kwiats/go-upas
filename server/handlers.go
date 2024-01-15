package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kwiats/go-upas/models"
	"github.com/kwiats/go-upas/utils"
)

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetUser(w, r)
	case "POST":
		return s.handleCreateUser(w, r)
	case "PATCH":
		return s.handleUpdateUser(w, r)
	default:
		return fmt.Errorf("not allowed %s methods", r.Method)
	}
}

func (s *APIServer) handleUserWithId(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetUserById(w, r)
	case "DELETE":
		return s.handleDeleteUser(w, r)
	default:
		return fmt.Errorf("not allowed %s methods", r.Method)
	}
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginUserRequest := new(models.UserLogin)

	if err := json.NewDecoder(r.Body).Decode(loginUserRequest); err != nil {
		log.Printf("Username and Password is required. %v", err)
		return models.WriteJSONResponse(w, http.StatusBadRequest, models.APIError{Error: "Username and Password is required"})
	}

	token, err := s.store.SignInUser(loginUserRequest)
	if err != nil {
		return models.WriteJSONResponse(w, http.StatusFound, err.Error())
	}

	return models.WriteJSONResponse(w, http.StatusOK, token)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		log.Printf("Cannot get all users. %v", err)
		return models.WriteJSONResponse(w, http.StatusNotFound, make([]string, 0))
	}
	return models.WriteJSONResponse(w, http.StatusFound, users)
}
func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	if id, ok := mux.Vars(r)["id"]; ok {
		user, err := s.store.GetUserByID(id)
		if err != nil {
			return models.WriteJSONResponse(w, http.StatusNotFound, err)
		}
		return models.WriteJSONResponse(w, http.StatusFound, user)
	}
	return models.WriteJSONResponse(w, http.StatusFound, "Id not provided")
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserRequest := new(models.UserCreate)
	if err := json.NewDecoder(r.Body).Decode(createUserRequest); err != nil {
		return err
	}

	user, err := createUserRequest.CreateAccount()
	if err != nil {
		return models.WriteJSONResponse(w, http.StatusBadRequest, err)
	}
	if err := s.store.CreateUser(user); err != nil {
		return models.WriteJSONResponse(w, http.StatusBadRequest, "Cannot create account with this credentials")

	}
	return models.WriteJSONResponse(w, http.StatusFound, user)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	if id, ok := mux.Vars(r)["id"]; ok {
		err := s.store.DeleteUser(id)
		if err != nil {
			return models.WriteJSONResponse(w, http.StatusNotFound, err)
		}
		return models.WriteJSONResponse(w, http.StatusOK, "User deleted")
	}
	return models.WriteJSONResponse(w, http.StatusFound, "Id not provided")
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) SignInUser(loginUserRequest *models.UserLogin) (*utils.TokenResponse, error) {
	user, err := s.store.GetUserCredentials(loginUserRequest)
	if err != nil {
		log.Printf("Cannot find user with this credentials. %v", err)
		return nil, fmt.Errorf("cannot find user with this credentials %v", err)
	}
	if ok := utils.CheckPassword(loginUserRequest.Password, user.Password); !ok {
		log.Printf("Incorrect password.")
		return nil, fmt.Errorf("incorrect password")
	}
	return utils.CreateJWT(user.ID, user.Username)

}
