package database

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"session-auth/utils"
)

// ConnectDB initializes a MongoDB connection.
func ConnectDB() (*mongo.Client, error) {
	log := utils.InitLogger()
	uri := viper.GetString("database.mongodburi")
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connection is established
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Info("Connected to MongoDB!")

	client.Database("sessionAuth").Collection("users")
	return client, nil
}
