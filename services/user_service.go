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
func CreateUser(db *mongo.Client, user *models.CreateUser) (*models.AuthResponse, *models.SessionError) {
	collection := db.Database("sessionAuth").Collection("users")

	// Check if the user already exists
	filter := map[string]interface{}{
		"email": user.Email,
	}

	// FineOne will return a nil error if a document is found
	var existingUser models.CreateUser
	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &models.SessionError{
				Message:     "User already exists",
				Status:      http.StatusConflict,
				Description: "A user with the provided email already exists. Please try logging in.",
				Errors:      "",
			}
		} else {
			return nil, &models.SessionError{
				Message:     "Internal server error",
				Description: "An internal server error occurred. Please try again later.",
				Status:      http.StatusInternalServerError,
				Errors:      "",
			}
		}
	}

	user.CreatedAt = time.Now()

	// Hash the password
	hashedPassword, existing := HashPassword(user.Password)
	if existing != nil {
		return nil, &models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}
	user.Password = hashedPassword

	// Insert the user into the database
	_, existing = collection.InsertOne(context.TODO(), user)
	if existing != nil {
		return nil, &models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      existing.Error(),
		}
	}

	// Generate a token for the user on successful creation
	token, err := GenerateJwt(user.Email)
	if err != nil {
		return nil, &models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}

	return &models.AuthResponse{Token: token}, nil
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

	// Generate a token for the user on successful creation
	token, err := GenerateJwt(user.Email)
	if err != nil {
		return nil, &models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}

	return &models.AuthResponse{Token: token}, nil
}
