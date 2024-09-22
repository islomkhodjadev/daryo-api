package controllers

import (
	"daryo-api/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	user, has := c.Get("user")

	if has {
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": "Authorization required"})
	}
	user = user.(*database.User)
	userID, has := c.Get("userID")
	if !has {
		c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": "User ID not found"})
		return
	}
	userIDStr := userID.(string)
	parsedUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = database.DeleteUser(parsedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

}
