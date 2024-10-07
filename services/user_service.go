package services

import (
	"context"
	"session-auth/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *mongo.Client, user *models.User) (interface{}, error) {
	collection := db.Database("sessionAuth").Collection("users")

	user.CreatedAt = time.Now()

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// Insert the user into the database
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
