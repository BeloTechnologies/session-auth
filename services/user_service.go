package services

import (
	"context"
	"session-auth/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *mongo.Client, user *models.CreateUser) (models.AuthResponse, error) {
	collection := db.Database("sessionAuth").Collection("users")

	user.CreatedAt = time.Now()

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return models.AuthResponse{}, err
	}
	user.Password = hashedPassword

	// Insert the user into the database
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Generate a token for the user on successful creation

	return models.AuthResponse{Token: "token"}, nil
}

func LoginUser(db *mongo.Client, user *models.LoginUser) (models.AuthResponse, error) {
	collection := db.Database("sessionAuth").Collection("users")

	// Find the user in the database
	filter := map[string]interface{}{
		"email": user.Email,
	}

	var result models.CreateUser
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Compare the stored password hash with the input password
	if !ComparePasswords(result.Password, user.Password) {
		return models.AuthResponse{}, nil
	}

	//// Generate a JWT token
	//token, err := GenerateToken(result.Username)
	//if err != nil {
	//	return "", err
	//}

	return models.AuthResponse{Token: "token"}, nil
}
