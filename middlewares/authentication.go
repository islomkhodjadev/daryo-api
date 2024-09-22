package middlewares

import (
	"daryo-api/database"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware is a middleware that checks if the token is valid
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Check if the token starts with "Bearer"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Verify the token
		userID, err := database.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		user, err := database.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
			c.Abort()
			return
		}
		// Set the user ID in the context
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	}
}
