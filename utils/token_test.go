package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	t.Setenv("TOKEN_HOUR_LIFESPAN", "12")
	t.Setenv("API_SECRET", "37")
	token, _ := GenerateToken(0, 100)
	assert.IsType(t, "", token)

	claims, _ := ValidToken(token)
	log.Println(claims)
	assert.Equal(t, float64(100), (*claims)["iss"])
	assert.Equal(t, float64(0), (*claims)["userId"])
}
