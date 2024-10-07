package controllers

import (
	"net/http"
	"session-auth/models"
	"session-auth/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateUser handles user creation.
func CreateUser(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.CreateUser

		// Validate input
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
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
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Message:     "Internal server error",
				Errors:      err.Error(),
				Status:      http.StatusInternalServerError,
				Description: "An internal server error occurred. Please try again later.",
			})
			return
		}

		c.JSON(http.StatusCreated, models.SuccessResponse{
			Message: "User created successfully",
			Status:  http.StatusCreated,
			Data:    result,
		})
	}
}

func LoginUser(db *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.LoginUser

		// Validate input
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the service to create a user
		token, err := services.LoginUser(db, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
