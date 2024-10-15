package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
)

// InitConfig function to initialize viper configuration
func InitConfig() {
	profile := "dev" // Dynamically change this based on your environment
	viper.SetConfigName(fmt.Sprintf("config.%s", profile))
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading configs file, %s", err)
	}

	LoadEnv()
}

// LoadEnv function to load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
}
