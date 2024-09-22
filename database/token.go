package database

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key for signing JWT tokens
var jwtKey = []byte("my_secret_key")

// Claims represents the payload data inside the JWT token
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for the user without expiration
func GenerateToken(userID int) (string, error) {
	// Create the JWT claims, including the user ID and the current timestamp
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(), // Add current timestamp to make the token unique
		},
	}

	// Create the token using HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// StoreToken stores a JWT token for a user in the database
func StoreToken(userID int, token string) error {
	query := `
	INSERT INTO user_tokens (user_id, token, created_at)
	VALUES ($1, $2, NOW())
	ON CONFLICT (user_id) DO UPDATE 
	SET token = $2, created_at = NOW();
	
	`
	_, err := DB.Exec(query, userID, token)
	if err != nil {
		fmt.Printf("error storing token: %v", err)
		return fmt.Errorf("error storing token: %v", err)
	}
	return nil
}

// VerifyToken verifies the JWT token and extracts user ID
func VerifyToken(tokenString string) (int, error) {
	claims := &Claims{}

	// Parse the token and validate it
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

// GetTokenByUserID retrieves the token for a specific user
func GetTokenByUserID(userID int) (string, error) {
	var token string
	query := `SELECT token FROM user_tokens WHERE user_id = $1;`
	err := DB.QueryRow(query, userID).Scan(&token)
	if err != nil {
		return "", fmt.Errorf("error fetching token: %v", err)
	}
	return token, nil
}

// UpdateToken generates a new token and updates the existing token for a user
func UpdateToken(userID int) (string, error) {
	// Generate a new token
	newToken, err := GenerateToken(userID)
	if err != nil {
		return "", fmt.Errorf("error generating new token: %v", err)
	}
	fmt.Print(newToken)
	// Store the new token in the database
	err = StoreToken(userID, newToken)
	if err != nil {
		return "", fmt.Errorf("error updating token: %v", err)
	}

	// Fetch the newly stored token from the database to ensure you're getting the latest token
	updatedToken, err := GetTokenByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("error fetching updated token: %v", err)
	}

	return updatedToken, nil
}
