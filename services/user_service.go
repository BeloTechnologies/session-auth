package services

import (
	"context"
	"session-auth/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *mongo.Client, user *models.CreateUser) (interface{}, error) {
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

func LoginUser(db *mongo.Client, user *models.LoginUser) (string, error) {
	collection := db.Database("sessionAuth").Collection("users")

	// Find the user in the database
	filter := map[string]interface{}{
		"email": user.Email,
	}

	var result models.CreateUser
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return "No user associated with this email: ", err
	}

	// Compare the stored password hash with the input password
	if !ComparePasswords(result.Password, user.Password) {
		return "Passwords do not match", nil
	}

	//// Generate a JWT token
	//token, err := GenerateToken(result.Username)
	//if err != nil {
	//	return "", err
	//}

	return "token", nil
}
