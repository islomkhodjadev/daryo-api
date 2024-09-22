package controllers

import (
	"daryo-api/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginHandler handles user login and token generation
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate input
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Authenticate the user (validate credentials)
	user, err := database.AuthenticateUser(username, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check if a valid token already exists
	token, err := database.GetTokenByUserID(user.ID)
	if err != nil {
		// No valid token, or error fetching token, so generate a new one
		token, err = database.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Store the new token in the database
		err = database.StoreToken(user.ID, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
			return
		}
	}

	// Respond with the token (either new or existing)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
