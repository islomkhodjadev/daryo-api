package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode,
	)
	DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Error while opening the db", err)
	}

	if err = DB.Ping(); err != nil {
		DB.Close()
		log.Fatal("Error db not opened", err)
	}
	fmt.Println("Connected to PostgreSQL database successfully!")

	err = InitializeDB()
	if err != nil {
		log.Fatal("Error while opening the db", err)
	}

}

func InitializeDB() error {
	query := `

	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,   
		email VARCHAR(100) UNIQUE NOT NULL,
		is_super_user BOOLEAN NOT NULL DEFAULT false,
		password VARCHAR(255) NOT NULL,         
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   
	);
	
	CREATE TABLE IF NOT EXISTS user_tokens (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,  -- Foreign key to users table
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (user_id)  -- Ensures each user can have only one token
	);
	
	

	CREATE TABLE IF NOT EXISTS clients (
		id SERIAL PRIMARY KEY,               -- Internal unique ID
		user_id INTEGER NOT NULL,            -- External service user ID
		name VARCHAR(100) NOT NULL,          -- Client's name
		email VARCHAR(100) UNIQUE NOT NULL,  -- Client's email
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS conversations (
		id SERIAL PRIMARY KEY,                -- Internal unique ID
		client_id INTEGER UNIQUE REFERENCES clients(id) ON DELETE CASCADE,  -- Ensure one conversation per client
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,                -- Internal unique ID
		conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,  -- Foreign key to conversations
		message TEXT NOT NULL,                -- Message text
		is_user_message BOOLEAN NOT NULL,     -- True if it's a user message
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);	
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}
	fmt.Println("Tables for users (businesses), clients, conversations, and messages are ready!")
	return nil
}
