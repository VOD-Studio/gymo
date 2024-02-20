package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
	"rua.plus/gymo/utils"
)

type User struct {
	Db *gorm.DB
}

type UserResponse struct {
	Message  string `json:"message"`
	Status   string `json:"status"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserQuery struct {
	Email string `form:"email" binding:"required"`
}

func (user User) GetUser(c *gin.Context) {
	var userInfo UserQuery
	if err := c.ShouldBindQuery(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userInfo.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is empty"})
		return
	}

	type result struct {
		ID       uint   `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	u := &result{}
	res := user.Db.Model(&models.User{}).Find(&u, "email = ?", userInfo.Email)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   u,
	})
}

type UserJson struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required"`
}

func (user User) AddUser(c *gin.Context) {
	var userInfo UserJson
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &models.User{
		Username: userInfo.Username,
		Password: userInfo.Password,
		Email:    userInfo.Email,
	}
	res := user.Db.Model(&models.User{}).Where("email = ?", u.Email).FirstOrCreate(&u)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exist"})
		return
	}

	type data struct {
		ID       uint   `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	resData := &data{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   resData,
	})
}

func (user User) ModifyUser(c *gin.Context) {
	var claims *jwt.MapClaims
	claim, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parse token failed"})
		return
	}
	claims = claim.(*jwt.MapClaims)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"userId": (*claims)["userId"],
	})

}

type UserLogin struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user User) Login(c *gin.Context) {
	var userInfo UserLogin
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// query the user
	u := &models.User{}
	dbResult := user.Db.Model(&models.User{}).Find(&u, "email = ?", userInfo.Email)
	if dbResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": dbResult.Error.Error()})
		return
	}
	if dbResult.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user not exist",
		})
		return
	}

	// check the password
	if err := models.CheckPasswordHash(userInfo.Password, u.Password); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "password not correct",
		})
		return
	}

	// generate token
	lastLogin := time.Now().Unix()
	token, err := utils.GenerateToken(int(u.ID), lastLogin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	type result struct {
		ID       uint   `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	res := &result{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
	}

	// update last login
	u.LastLogin = lastLogin
	user.Db.Save(&u)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   res,
	})
}
