package models

import "time"

type Conversation struct {
	ID        int       `json:"id"`         // Unique ID of the conversation
	ClientID  int       `json:"client_id"`  // References the client this conversation belongs to
	CreatedAt time.Time `json:"created_at"` // When the conversation started
}
