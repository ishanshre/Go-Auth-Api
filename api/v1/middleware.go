package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ishanshre/go-auth-api/api/v1/models"
	"github.com/ishanshre/go-auth-api/internals/pkg/utils"
)

func makeHttpHandler(f ApiFunc) http.HandlerFunc {
	/*
		Middleware that returns http HandlerFunc
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	/*
		Middleware that write the response to the client
	*/
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getId(r *http.Request) (int, error) {
	/*
		Middlere that parse the id in url params into interger
	*/
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("error in parsing the id")
	}
	return id, nil
}

func scanUsers(rows *sql.Rows) (*models.User, error) {
	/*
		A middleware that scan the values in record and store in another variable
	*/
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

func jwtAuthHandler(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	/*
		Authentication middleware that handles the authorizations
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling auth middleware")
		userId, err := utils.ExtractTokenMetaData(r)
		if err != nil {
			log.Println(err)
			permissionDenied(w)
			return
		}
		account, err := s.GetUsersById(userId.ID)
		if err != nil {
			log.Println(err)
			permissionDenied(w)
			return
		}
		if err := utils.VerifyUser(account.ID, r); err != nil {
			log.Println(err)
			permissionDenied(w)
			return
		}
		paramId, err := getId(r)
		if err != nil {
			invalidParams(w)
			return
		}
		if paramId != account.ID {
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

func jwtAuthAdminHandler(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	/*
		Authorization Handler for admin
	*/
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Calling auth middleware")
		userId, err := utils.ExtractTokenMetaData(r)
		if err != nil {
			log.Println(err)
			permissionDenied(w)
			return
		}
		account, err := s.GetUsersById(userId.ID)
		if err != nil {
			log.Println(err)
			permissionDenied(w)
			return
		}
		if err := utils.VerifyUser(account.ID, r); err != nil {
			log.Println(err)
			permissionDenied(w)
			return
		}
		if account.IsAdmin == "false" {
			log.Println("permission denied, not a admin user")
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}
