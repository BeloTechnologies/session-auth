package main

import (
	"log"
	"session-auth/config"
	"session-auth/database"
	"session-auth/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	r := gin.Default()

	routes.UserRoutes(r, db)

	log.Println("Starting server on :8080")

	r.Run(":8080")
}
