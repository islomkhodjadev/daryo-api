package models

import "time"

type Message struct {
	ID             int       `json:"id"`              // Unique ID of the message
	ConversationID int       `json:"conversation_id"` // References the conversation this message belongs to
	Message        string    `json:"message"`         // The actual message text
	IsUserMessage  bool      `json:"is_user_message"` // True if it's from the client, false if it's from the AI
	CreatedAt      time.Time `json:"created_at"`      // When the message was sent
}
