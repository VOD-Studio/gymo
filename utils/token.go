package utils

import (
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user_id int) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"authorized": true,
		"user_id":    user_id,
		"exp":        time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}
