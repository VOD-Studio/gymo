package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
	"rua.plus/gymo/utils"
)

// 验证 Token 并将其添加到上下文中
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		tokenArray := strings.Split(token, " ")

		if tokenArray[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}

		claims, err := utils.ValidToken(tokenArray[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// 获取 token 并查询到对应到用户
// 添加到上下文中
func TokenTimeAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *jwt.MapClaims
		claim, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "parse token failed"},
			)
			return
		}
		claims = claim.(*jwt.MapClaims)

		user := &models.User{}
		res := db.Model(user).Find(user, "id = ?", (*claims)["userId"])
		if res.Error != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": res.Error.Error()},
			)
			return
		}
		if res.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
			})
			return
		}
		if user.LastLogin != int64((*claims)["iss"].(float64)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "token expired",
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}

}
