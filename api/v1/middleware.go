package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ishanshre/go-auth-api/api/v1/models"
)

func makeHttpHandler(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("error in parsing the id")
	}
	return id, nil
}

func scanUsers(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)
	return user, err
}

func scanUser1(rows *sql.Rows) (*models.UserNhash, error) {
	user := new(models.UserNhash)
	err := rows.Scan(
		&user.ID,
		&user.Password,
	)
	return user, err
}
