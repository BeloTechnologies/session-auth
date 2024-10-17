package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
	"session-auth/configs"
	"session-auth/database"
	"session-auth/routes"
	"session-auth/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Init logger
	log := utils.InitLogger()
	log.Info("Initializing server...")

	log.Info("Loading environment variables and configs...")
	configs.LoadEnv()
	utils.InitConfig()

	serverPort := viper.GetInt("server.port")

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
	routes.AuthRoutes(r)

	log.Info(fmt.Sprintf("Starting server on :%d", serverPort))

	e := r.Run(fmt.Sprintf(":%d", serverPort))
	if e != nil {
		return
	}
}
