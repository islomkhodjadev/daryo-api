package models

import "time"

// Client struct represents a client associated with a business (user)
type User struct {
	ID        int       `json:"id"`         // Unique ID of the client
	Name      string    `json:"name"`       // Client's name
	Email     string    `json:"email"`      // Client's email address
	CreatedAt time.Time `json:"created_at"` // When the client was created
}
