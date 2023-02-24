package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func ExtractToken(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	tokenString := strings.Split(bearerToken, " ")
	if len(tokenString) == 2 {
		return tokenString[1], nil
	}
	return "", fmt.Errorf("invalid token")
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := ExtractToken(r)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) (jwt.MapClaims, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, err
	}
	return claims, nil
}

func ExtractTokenMetaData(r *http.Request) (*models.TokenMetaData, error) {
	claims, err := TokenValid(r)
	if err != nil {
		return nil, err
	}
	userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 0)
	if err != nil {
		return nil, err
	}
	return &models.TokenMetaData{ID: int(userId)}, nil
}

func VerifyUser(id int, r *http.Request) error {
	claims, err := TokenValid(r)
	if err != nil {
		return err
	}
	if int64(id) != int64(claims["user_id"].(float64)) {
		return fmt.Errorf("permission denied")
	}
	return nil
}
