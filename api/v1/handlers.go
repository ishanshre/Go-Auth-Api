package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ishanshre/go-auth-api/api/v1/models"
)

func (s *ApiServer) handleCreateAndGetUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUsers(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateNewUser(w, r)
	}
	return fmt.Errorf("%s method not allowed", r.Method)
}

func (s *ApiServer) handleUsersById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUserById(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateUserById(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUserById(w, r)
	}
	return fmt.Errorf("%s method not allowed", r.Method)
}

func (s *ApiServer) handleCreateNewUser(w http.ResponseWriter, r *http.Request) error {
	req := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	if err := validateCreateUser(req); err != nil {
		return err
	}
	user, err := models.NewUser(
		req.FirstName,
		req.LastName,
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}
	if err := s.store.CreateUser(user); err != nil {
		return err
	}
	return writeJSON(w, http.StatusCreated, ApiSuccess{Success: "account created successfully"})
}

func (s *ApiServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	account, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	user, err := s.store.GetUsersById(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, user)
}

func (s *ApiServer) handleUpdateUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	req := new(models.UpdateUser)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	if err := s.store.UpdateUserById(id, req); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, ApiSuccess{Success: "user updated"})
}

func (s *ApiServer) handleDeleteUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteUserById(id); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, ApiSuccess{Success: "user deleted"})
}
