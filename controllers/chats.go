package controllers

import (
	"daryo-api/database"
	"daryo-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleClientMessage(c *gin.Context) {
	userIDStr := c.PostForm("user_id")
	name := c.PostForm("name")
	email := c.PostForm("email")
	messageText := c.PostForm("message")

	if userIDStr == "" || email == "" || name == "" || messageText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, email, name, and message are required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	// Step 1: Check if the client exists by email
	client, err := database.GetClientByEmail(email)
	if err != nil {
		client, err = database.CreateClient(userID, name, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new client"})
			return
		}
	}

	// Step 2: Check if the client has an existing active conversation
	conversation, err := database.GetLatestConversation(client.ID)
	if err != nil {
		// No conversation exists, so create a new one
		conversation, err = database.CreateConversation(client.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new conversation"})
			return
		}
	}

	// Call AddMessage and ignore the newMessage value, only check for err
	if _, err := database.AddMessage(conversation.ID, messageText, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add message"})
		return
	}

	allMessage, err := database.GetAllMessagesAsString(conversation.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	aiResponse, err := utils.Gpt(allMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ai returned wrong answer"})
		return
	}

	database.AddMessage(conversation.ID, aiResponse, false)

	// Return success with conversation and message details
	c.JSON(http.StatusOK, gin.H{
		"client":       client,
		"conversation": conversation,
		"message":      aiResponse,
	})
}
