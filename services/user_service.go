package services

import (
	"context"
	"errors"
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/BeloTechnologies/session-core/core_models/user_models"
	"net/http"
	"session-auth/models"
	"session-auth/utils"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *mongo.Client, user *models.CreateUser) (*models.AuthResponse, *core_models.SessionError) {
	log := utils.InitLogger()
	collection := db.Database("sessionAuth").Collection("users")

	// Check if the user already exists by email
	filter := map[string]interface{}{
		"email": user.Email,
	}
	// TODO: Check if the user exists via username and phone too!

	// FineOne will return a nil error if a document is found
	var existingUser models.CreateUser
	err := collection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		log.Errorf("User already exists: %v", err)
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &core_models.SessionError{
				Message:     "User already exists",
				Status:      http.StatusConflict,
				Description: "A user with the provided email already exists. Please try logging in.",
				Errors:      "",
			}
		} else {
			return nil, &core_models.SessionError{
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
		log.Errorf("Error hashing password: %v", existing)
		return nil, &core_models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}
	user.Password = hashedPassword

	// Call the session-user service to create an entry in relational database
	createUserRowResult, proxyErr := CreateUserEntry(user_models.CreateUserRow{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
	})
	if proxyErr != nil {
		log.Errorf("Error creating user row: %v", proxyErr)
		return nil, &core_models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}

	// Set the psqlID in the user struct
	user.PsqlID = createUserRowResult.ID

	// Insert the user into the database
	_, existing = collection.InsertOne(context.TODO(), user)
	if existing != nil {
		return nil, &core_models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      existing.Error(),
		}
	}

	// Generate a token for the user on successful creation
	token, err := GenerateJwt(strconv.Itoa(user.PsqlID))
	if err != nil {
		return nil, &core_models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}

	createUserRowResult.ID = -1

	return &models.AuthResponse{
		Token:    token,
		UserData: createUserRowResult,
	}, nil
}

func LoginUser(db *mongo.Client, user *models.LoginUser) (*models.AuthResponse, *core_models.SessionError) {
	log := utils.InitLogger()
	collection := db.Database("sessionAuth").Collection("users")

	// Find the user in the database
	filter := map[string]interface{}{
		"email": user.Email,
	}

	var result models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Errorf("Error finding user: %v", err)
		return nil, &core_models.SessionError{
			Message:     "User not found",
			Status:      http.StatusNotFound,
			Description: "The user with the provided email does not exist.",
			Errors:      err.Error(),
		}
	}

	log.Infof("User found: %v", result)

	// Compare the stored password hash with the input password
	if !ComparePasswords(result.Password, user.Password) {
		return nil, &core_models.SessionError{
			Message:     "Invalid credentials",
			Status:      http.StatusBadRequest,
			Description: "The email or password provided is incorrect.",
		}
	}

	// Generate a token for the user on successful creation
	token, err := GenerateJwt(user.Email)
	if err != nil {
		return nil, &core_models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}

	// Get User details from the relational database for fast access
	userRow, proxyErr := GetUser(result.PsqlID, token)
	if proxyErr != nil {
		log.Errorf("Error getting user row: %v", proxyErr)
		return nil, &core_models.SessionError{
			Message:     "Internal server error",
			Description: "An internal server error occurred. Please try again later.",
			Status:      http.StatusInternalServerError,
			Errors:      "",
		}
	}

	userRow.ID = -1

	return &models.AuthResponse{
		Token:    token,
		UserData: userRow,
	}, nil
}
