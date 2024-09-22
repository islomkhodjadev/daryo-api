package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Client struct
type Client struct {
	ID        int
	UserID    int
	Name      string
	Email     string
	CreatedAt time.Time
}

func CreateClient(userID int, name, email string) (*Client, error) {
	query := `
		INSERT INTO clients (user_id, name, email, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, user_id, name, email, created_at;
	`
	var client Client
	err := DB.QueryRow(query, userID, name, email).Scan(&client.ID, &client.UserID, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	return &client, nil
}

func GetAllClients() ([]Client, error) {
	query := `SELECT id, user_id, name, email, created_at FROM clients;`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching clients: %v", err)
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.UserID, &client.Name, &client.Email, &client.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning client: %v", err)
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func GetClientByID(clientID int) (*Client, error) {
	query := `SELECT id, user_id, name, email, created_at FROM clients WHERE id = $1;`
	var client Client
	err := DB.QueryRow(query, clientID).Scan(&client.ID, &client.UserID, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("error fetching client: %v", err)
	}
	return &client, nil
}

func GetClientByEmail(clientEmail string) (*Client, error) {
	query := `SELECT id, user_id, name, email, created_at FROM clients WHERE email = $1;`
	var client Client
	err := DB.QueryRow(query, clientEmail).Scan(&client.ID, &client.UserID, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("error fetching client: %v", err)
	}
	return &client, nil
}

// Update a client
func UpdateClient(clientID int, name, email string) (*Client, error) {
	query := `
		UPDATE clients SET name = $1, email = $2 WHERE id = $3
		RETURNING id, user_id, name, email, created_at;
	`
	var client Client
	err := DB.QueryRow(query, name, email, clientID).Scan(&client.ID, &client.UserID, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating client: %v", err)
	}
	return &client, nil
}

// Delete a client
func DeleteClient(clientID int) error {
	query := `DELETE FROM clients WHERE id = $1;`
	_, err := DB.Exec(query, clientID)
	if err != nil {
		return fmt.Errorf("error deleting client: %v", err)
	}
	return nil
}
