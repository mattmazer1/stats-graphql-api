package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte

func GetSecretKey() ([]byte, error) {
	if secretKey != nil {
		return secretKey, nil
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("JWT_SECRET environment variable not set")
		os.Exit(1)
	}

	secretKey = []byte(secret)
	return secretKey, nil
}

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secret, err := GetSecretKey()
	if err != nil {
		return "", fmt.Errorf("could not get private key %w", err)
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {

		return "", fmt.Errorf("error in Generating key %w", err)
	}

	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username it claims
func ParseToken(tokenStr string) (string, error) {
	secret, err := GetSecretKey()
	if err != nil {
		return "", fmt.Errorf("could not get private key %w", err)
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("could not parse jwt %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
