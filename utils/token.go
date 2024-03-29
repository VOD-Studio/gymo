package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId int, lastLogin int64) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iss":    lastLogin,
		"exp":    time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func ValidToken(tokenString string) (*jwt.MapClaims, error) {
	clamis := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &clamis, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return &clamis, nil
}
