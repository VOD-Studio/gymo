package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"rua.plus/gymo/utils"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		tokenArray := strings.Split(token, " ")

		if tokenArray[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
		}

		claims, err := utils.ValidToken(tokenArray[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
		}
		c.Set("claims", claims)
		c.Next()
	}
}
