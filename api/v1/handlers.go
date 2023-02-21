package v1

import (
	"encoding/json"
	"net/http"

	"github.com/ishanshre/go-auth-api/api/v1/models"
)

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
