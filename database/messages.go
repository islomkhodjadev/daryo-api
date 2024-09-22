package database

import (
	"fmt"
	"time"
)

// Message struct
type Message struct {
	ID             int
	ConversationID int
	Message        string
	IsUserMessage  bool
	CreatedAt      time.Time
}

// Add a message to a conversation
func AddMessage(conversationID int, message string, isUserMessage bool) (*Message, error) {
	query := `
		INSERT INTO messages (conversation_id, message, is_user_message, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, conversation_id, message, is_user_message, created_at;
	`
	var msg Message
	err := DB.QueryRow(query, conversationID, message, isUserMessage).Scan(&msg.ID, &msg.ConversationID, &msg.Message, &msg.IsUserMessage, &msg.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error adding message: %v", err)
	}
	return &msg, nil
}

// Get all messages in a conversation
func GetAllMessages(conversationID int) ([]Message, error) {
	query := `SELECT id, conversation_id, message, is_user_message, created_at FROM messages WHERE conversation_id = $1;`
	rows, err := DB.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %v", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Message, &msg.IsUserMessage, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning message: %v", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func GetAllMessagesAsString(conversationID int) (string, error) {
	query := `SELECT id, conversation_id, message, is_user_message, created_at FROM messages WHERE conversation_id = $1;`
	rows, err := DB.Query(query, conversationID)
	if err != nil {
		return "", fmt.Errorf("error fetching messages: %v", err)
	}
	defer rows.Close()

	var conversationStr string
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Message, &msg.IsUserMessage, &msg.CreatedAt); err != nil {
			return "", fmt.Errorf("error scanning message: %v", err)
		}
		// Format message based on whether it's a user message or AI message
		if msg.IsUserMessage {
			conversationStr += fmt.Sprintf("User: %s\n", msg.Message)
		} else {
			conversationStr += fmt.Sprintf("AI: %s\n", msg.Message)
		}
	}
	return conversationStr, nil
}

// Delete a message
func DeleteMessage(messageID int) error {
	query := `DELETE FROM messages WHERE id = $1;`
	_, err := DB.Exec(query, messageID)
	if err != nil {
		return fmt.Errorf("error deleting message: %v", err)
	}
	return nil
}
