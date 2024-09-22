package controllers

import (
	"daryo-api/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	fmt.Print("salom")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	is_super_user := c.PostForm("is_super_user")

	fmt.Print(username, email, password)
	// Validate input fields
	if username == "" || email == "" || password == "" || is_super_user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
		return
	}
	isSuperUserBool := false
	if is_super_user == "true" {
		isSuperUserBool = true
	}

	existingUser, _ := database.GetUserByEmail(email)
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already registered"})
		return
	}

	user, err := database.CreateUser(username, email, password, isSuperUserBool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Respond with the user data (excluding password)
	c.JSON(http.StatusOK, gin.H{
		"message":  "User registered successfully",
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
