package controllers

import (
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"session-auth/models"
	"session-auth/services"
	"session-auth/utils"
)

// CreateUser handles user creation.
func CreateUser(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := utils.InitLogger() // Initialize and get the logger

		var user models.CreateUser
		// Validate input
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Errorf("Invalid input: %v", err)
			c.JSON(http.StatusBadRequest, core_models.ErrorResponse{
				Message:     "Invalid input",
				Errors:      err.Error(),
				Status:      http.StatusBadRequest,
				Description: "The input provided is invalid. Please check the input and try again.",
			})
			return
		}

		// Call the service to create a user
		result, err := services.CreateUser(db, &user)
		if err != nil {
			log.Errorf("Error creating user: %v", err)
			c.JSON(err.Status, core_models.ErrorResponse{
				Message:     err.Message,
				Errors:      err.Errors,
				Status:      err.Status,
				Description: err.Description,
			})
			return
		}

		log.Info("User created successfully")
		c.JSON(http.StatusCreated, core_models.SuccessResponse{
			Message: "User created successfully",
			Status:  http.StatusCreated,
			Data:    result,
		})
	}
}

// LoginUser handles user login.
func LoginUser(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := utils.InitLogger() // Initialize and get the logger

		var user models.LoginUser
		// Validate input
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Errorf("Invalid input: %v", err)
			c.JSON(http.StatusBadRequest, core_models.ErrorResponse{
				Message:     "Invalid input",
				Errors:      err.Error(),
				Status:      http.StatusBadRequest,
				Description: "The input provided is invalid. Please check the input and try again.",
			})
			return
		}

		// Call the service to login the user
		result, err := services.LoginUser(db, &user)
		if err != nil {
			log.Errorf("Error logging in user: %v", err)
			c.JSON(http.StatusInternalServerError, core_models.ErrorResponse{
				Message:     err.Message,
				Errors:      err.Errors,
				Status:      err.Status,
				Description: err.Description,
			})
			return
		}

		log.Info("User logged in successfully")
		c.JSON(http.StatusOK, core_models.SuccessResponse{
			Message: "User logged in successfully",
			Status:  http.StatusOK,
			Data:    result,
		})
	}
}
