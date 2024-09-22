package main

import (
	"daryo-api/database"
	"daryo-api/routers"
	"log"
)

func main() {
	// Step 1: Initialize Database Connection
	database.Connection() // Connect to the database and create tables if they don't exist
	defer database.DB.Close()

	// Step 2: Initialize Routes using Gin
	r := routers.InitializeRoutes()

	// Step 3: Start the Gin HTTP server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
