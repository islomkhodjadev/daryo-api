package database

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents a user in the system
type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Is_super_user bool      `json:"is_super_user"`
	CreatedAt     time.Time `json:"created_at"`
}

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the provided password matches the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(username, email, password string, is_super_user bool) (*User, error) {
	// Hash the password before storing
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	query := `
        INSERT INTO users (username, email, password, is_super_user)
        VALUES ($1, $2, $3, $4)
        RETURNING id, username, email, is_super_user created_at;
    `
	var user User

	// Run the query and scan the result into the User struct
	err = DB.QueryRow(query, username, email, hashedPassword, is_super_user).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		// Print the error to understand the issue
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID (without exposing the password)
func GetUserByID(userID int) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1;`
	var user User
	err := DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE email = $1;`
	var user User
	err := DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}
	return &user, nil
}

// AuthenticateUser checks if the provided username and password are correct
func AuthenticateUser(username, password string) (*User, error) {
	var user User
	var hashedPassword string

	query := `SELECT id, username, email, password, created_at FROM users WHERE username = $1`
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &hashedPassword, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	// Compare the stored hashed password with the provided password
	if !CheckPasswordHash(password, hashedPassword) {
		return nil, fmt.Errorf("invalid password")
	}

	// Return the user without the password
	user.Password = ""
	return &user, nil
}

// UpdateUser updates a user's username, email, or password
func UpdateUser(userID int, username, email, password string) (*User, error) {
	var hashedPassword string
	if password != "" {
		var err error
		hashedPassword, err = HashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("error hashing password: %v", err)
		}
	}

	query := `
		UPDATE users SET username = $1, email = $2, password = COALESCE($3, password)
		WHERE id = $4
		RETURNING id, username, email, created_at;
	`
	var user User
	err := DB.QueryRow(query, username, email, hashedPassword, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}
	return &user, nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
