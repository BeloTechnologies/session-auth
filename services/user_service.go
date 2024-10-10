package services

import (
	"context"
	"errors"
	"net/http"
	"session-auth/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *mongo.Client, user *models.CreateUser) (*models.AuthResponse, *models.ErrorResponse) {
	collection := db.Database("sessionAuth").Collection("users")

	// Check if the user already exists
	filter := map[string]interface{}{
		"email": user.Email,
	}

	var existingUser models.CreateUser
	existing := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if !errors.Is(existing, mongo.ErrNoDocuments) {
		return nil, &models.ErrorResponse{
			Message:     "User already exists",
			Description: "A user with the provided email already exists. Please try logging in.",
			Status:      http.StatusConflict,
			Errors:      "User already exists",
		}
	}

	user.CreatedAt = time.Now()

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, &models.ErrorResponse{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      err.Error(),
		}
	}
	user.Password = hashedPassword

	// Insert the user into the database
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, &models.ErrorResponse{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      err.Error(),
		}
	}

	// Generate a token for the user on successful creation

	return &models.AuthResponse{Token: "token"}, nil
}

func LoginUser(db *mongo.Client, user *models.LoginUser) (*models.AuthResponse, *models.SessionError) {
	collection := db.Database("sessionAuth").Collection("users")

	// Find the user in the database
	filter := map[string]interface{}{
		"email": user.Email,
	}

	var result models.CreateUser
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, &models.SessionError{
			Message:     "User not found",
			Status:      http.StatusNotFound,
			Description: "The user with the provided email does not exist.",
			Errors:      err.Error(),
		}
	}

	// Compare the stored password hash with the input password
	if !ComparePasswords(result.Password, user.Password) {
		return nil, &models.SessionError{
			Message:     "Invalid credentials",
			Status:      http.StatusBadRequest,
			Description: "The email or password provided is incorrect.",
		}
	}

	return &models.AuthResponse{Token: "token"}, nil
}
