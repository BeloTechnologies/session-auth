package main

import (
	"github.com/gin-contrib/cors"
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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"*"}

	r := gin.Default()

	r.Use(cors.New(config))

	routes.UserRoutes(r, db)

	log.Println("Starting server on :8080")

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
