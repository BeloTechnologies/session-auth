package routes

import (
	"github.com/gin-gonic/gin"
	"session-auth/controllers"
)

func AuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/ping/", controllers.Ping())
	}
}
