package controllers

import (
	"daryo-api/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateTokenHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate username and password are provided
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Step 1: Authenticate the user (validate credentials)
	user, err := database.AuthenticateUser(username, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Step 2: Generate a new token for the authenticated user
	token, err := database.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Step 3: Store the token in the database
	err = database.StoreToken(user.ID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UpdateTokenHandler(c *gin.Context) {

	user, has := c.Get("user")
	if !has {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Assuming user is of type *User (a pointer to a User struct)
	userStruct, ok := user.(*database.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast user"})
		return
	}

	userID := userStruct.ID // Now you can access user.ID

	newToken, err := database.UpdateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"new_token": newToken})
}
