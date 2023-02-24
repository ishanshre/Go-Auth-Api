package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsAdmin   string    `json:"is_admin"`
	IsActive  string    `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}
type RegisterUser struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsAdmin   string    `json:"is_admin"`
	IsActive  string    `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}

type AdminUpdateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsAdmin   string `json:"is_admin"`
	IsActive  string `json:"is_active"`
}

type UserNhash struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type UpdateUser struct {
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           int    `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenMetaData struct {
	ID int `json:"id"`
}

func NewUser(firstName, lastName, username, email, password string) (*User, error) {

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().Local().UTC(),
		UpdatedAt: time.Now().Local().UTC(),
		LastLogin: time.Time{},
	}, nil
}

func RegisterNewUser(firstName, lastName, username, email, password string) (*RegisterUser, error) {

	return &RegisterUser{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().Local().UTC(),
		UpdatedAt: time.Now().Local().UTC(),
		LastLogin: time.Time{},
	}, nil
}
