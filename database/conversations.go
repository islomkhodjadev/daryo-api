package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Conversation struct
type Conversation struct {
	ID        int
	ClientID  int
	CreatedAt time.Time
}

// GetLatestConversation retrieves the single conversation for a client
func GetLatestConversation(clientID int) (*Conversation, error) {
	query := `
		SELECT id, client_id, created_at
		FROM conversations
		WHERE client_id = $1
		LIMIT 1;
	`

	var conversation Conversation
	err := DB.QueryRow(query, clientID).Scan(&conversation.ID, &conversation.ClientID, &conversation.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no conversation found for this client")
		}
		return nil, fmt.Errorf("error fetching conversation: %v", err)
	}

	return &conversation, nil
}

// Create a new conversation for a client, or return the existing one
func CreateConversation(clientID int) (*Conversation, error) {
	// First, check if a conversation already exists for this client
	existingConversation, err := GetLatestConversation(clientID)
	if err == nil {
		// If conversation exists, return it
		return existingConversation, fmt.Errorf("conversation already exists for this client")
	}

	// If no existing conversation, create a new one
	query := `
		INSERT INTO conversations (client_id, created_at)
		VALUES ($1, NOW())
		RETURNING id, client_id, created_at;
	`
	var conversation Conversation
	err = DB.QueryRow(query, clientID).Scan(&conversation.ID, &conversation.ClientID, &conversation.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating conversation: %v", err)
	}
	return &conversation, nil
}

// Delete a conversation by its ID
func DeleteConversation(conversationID int) error {
	query := `DELETE FROM conversations WHERE id = $1;`
	_, err := DB.Exec(query, conversationID)
	if err != nil {
		return fmt.Errorf("error deleting conversation: %v", err)
	}
	return nil
}
