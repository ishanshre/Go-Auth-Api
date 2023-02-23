package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ishanshre/go-auth-api/api/v1/models"
)

func GenerateToken(id int) (*models.LoginResponse, error) {
	access_claims := &jwt.MapClaims{
		"ExpiresAt": jwt.NewNumericDate(time.Now().Add(time.Minute * 15)).Unix(),
		"IssuedAT":  jwt.NewNumericDate(time.Now()),
		"user_id":   id,
	}
	secret := os.Getenv("JWT_SECRET")
	ss := jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims)
	access_token, err := ss.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("error in generating access token %s", err)
	}
	refresh_claims := &jwt.MapClaims{
		"ExpiresAt": jwt.NewNumericDate(time.Now().Add(time.Hour)).Unix(),
		"IssuedAt":  jwt.NewNumericDate(time.Now()),
		"user_id":   id,
	}
	rss := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims)
	refresh_token, err := rss.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("error in generating refresh token %s", err)
	}
	return &models.LoginResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ID:           id,
	}, nil
}
