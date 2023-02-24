package utils

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	log.Println(password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	log.Println("hashed pass: ", string(hashedPassword))
	if err != nil {
		return "", fmt.Errorf("error in hashing the password")
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	log.Println(hashedPassword, candidatePassword)

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
