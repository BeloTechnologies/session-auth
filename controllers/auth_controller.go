package controllers

import (
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/BeloTechnologies/session-core/core_models/auth_models"
	"github.com/gin-gonic/gin"
	"net/http"
	"session-auth/services"
	"session-auth/utils"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := utils.InitLogger() // Initialize and get the logger

		var token auth_models.Token
		// Validate input
		if err := c.ShouldBindJSON(&token); err != nil {
			log.Errorf("Invalid input: %v", err)
			c.JSON(http.StatusBadRequest, core_models.ErrorResponse{
				Message:     "Invalid input",
				Errors:      err.Error(),
				Status:      http.StatusBadRequest,
				Description: "The input provided is invalid. Please check the input and try again.",
			})
			return
		}

		// Call the service to validate the token
		result, err := services.ValidateToken(token)
		if err != nil {
			log.Errorf("Error validating token: %v", err)
			c.JSON(err.Status, core_models.ErrorResponse{
				Message:     err.Message,
				Errors:      err.Errors,
				Status:      err.Status,
				Description: err.Description,
			})
			return
		}

		log.Info("Token validated successfully")
		c.JSON(http.StatusOK, core_models.SuccessResponse{
			Message: "Token is valid",
			Status:  http.StatusOK,
			Data:    result,
		})
	}
}
