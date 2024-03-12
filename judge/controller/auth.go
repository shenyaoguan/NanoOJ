package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xyz.xeonds/nano-oj/database/model"
	"xyz.xeonds/nano-oj/utils"
)

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"` //Email format authorized here
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := repo.GetUserByEmail(input.Email)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := repo.GetUserByUsername(user.Username); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	if err := repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	if user.ID == 1 { // if it is the first user, set it as admin
		user.AccountInfo.Permission = 0
		if err := repo.UpdateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func GetUser(c *gin.Context) {
	if userID := c.Param("user_id"); userID != "" {
		userID, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		user, err := repo.GetUserByUserID(uint32(userID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}
	users, err := repo.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if _, err := repo.GetUserByUsername(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	if err := repo.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := repo.GetUserByUsername(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}
	if err := repo.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
