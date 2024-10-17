package controllers

import (
	"github.com/gin-gonic/gin"
	"session-auth/utils"
)

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := utils.InitLogger() // Initialize and get the logger
		log.Info("Ping endpoint hit")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
