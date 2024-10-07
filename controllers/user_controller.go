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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the service to create a user
		insertedID, err := services.CreateUser(db, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user_id": insertedID})
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
